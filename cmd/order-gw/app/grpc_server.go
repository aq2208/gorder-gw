package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"gorder-gw/configs"
	"net"
	"os"
	"time"

	"gorder-gw/internal/controller/grpcapi"
	gwpb "gorder-gw/internal/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// RunGRPC starts the gRPC server and blocks until ctx is cancelled.
// It returns when the server shuts down gracefully.
func RunGRPC(ctx context.Context, cfg configs.Config, svc *grpcapi.OrderService) error {
	if cfg.GrpcServer.ListenAddr == "" {
		cfg.GrpcServer.ListenAddr = ":50051"
	}
	if cfg.GrpcServer.ShutdownGRace == 0 {
		cfg.GrpcServer.ShutdownGRace = 10 * time.Second
	}

	lis, err := net.Listen("tcp", cfg.GrpcServer.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	// Credentials
	var opts []grpc.ServerOption
	if cfg.GrpcServer.UseTLS {
		creds, err := buildServerCreds(cfg)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	// Interceptors (add logging / metrics here)
	//opts = append(opts,
	//	grpc.ChainUnaryInterceptor(unaryLogging()),
	//)

	grpcSrv := grpc.NewServer(opts...)

	// Health + reflection
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcSrv, healthSrv)
	reflection.Register(grpcSrv)

	// Register our service
	gwpb.RegisterOrderServiceServer(grpcSrv, svc)

	// Serve in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcSrv.Serve(lis)
	}()

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		// graceful stop with deadline
		done := make(chan struct{})
		go func() {
			grpcSrv.GracefulStop()
			close(done)
		}()
		select {
		case <-done:
			return nil
		case <-time.After(cfg.GrpcServer.ShutdownGRace):
			grpcSrv.Stop()
			return nil
		}
	case err := <-errCh:
		return err
	}
}

func buildServerCreds(cfg configs.Config) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(cfg.GrpcServer.CertFile, cfg.GrpcServer.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load cert/key: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Optional: enable mTLS if CAFile is provided (require client cert)
	if cfg.GrpcServer.CAFile != "" {
		caPEM, err := os.ReadFile(cfg.GrpcServer.CAFile)
		if err != nil {
			return nil, fmt.Errorf("read CA: %w", err)
		}
		cp := x509.NewCertPool()
		if ok := cp.AppendCertsFromPEM(caPEM); !ok {
			return nil, errors.New("bad CA file")
		}
		tlsCfg.ClientCAs = cp
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return credentials.NewTLS(tlsCfg), nil
}
