package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App struct {
		Name     string `koanf:"name"`
		HTTPAddr string `koanf:"http_addr"`
		LogLevel string `koanf:"log_level"`
	} `koanf:"app"`

	HTTP struct {
		ReadTimeout  time.Duration `koanf:"read_timeout"`
		WriteTimeout time.Duration `koanf:"write_timeout"`
		IdleTimeout  time.Duration `koanf:"idle_timeout"`
	} `koanf:"http"`

	GrpcServer struct {
		ListenAddr    string        `koanf:"listen_addr"`
		UseTLS        bool          `koanf:"use_tls"`
		CertFile      string        `koanf:"cert_file"`
		KeyFile       string        `koanf:"key_file"`
		CAFile        string        `koanf:"ca_file"`
		ShutdownGRace time.Duration `koanf:"shutdown_grace"`
	} `koanf:"grpc_server"`

	KafkaBroker struct {
		KafkaBrokers []string `koanf:"brokers"`
		KafkaTopic   string   `koanf:"topic"`
	} `koanf:"kafka"`
}

func Load(pathDir, envName string) (Config, error) {
	k := koanf.New(".")
	// 1) base
	//if err := k.Load(file.Provider(fmt.Sprintf("%s/base.yaml", pathDir)), yaml.Parser()); err != nil {
	//	return Config{}, fmt.Errorf("load base: %w", err)
	//}

	// 2) env override (dev/staging/prod). Optional: allow missing for local runs.
	_ = k.Load(file.Provider(fmt.Sprintf("%s/%s.yaml", pathDir, envName)), yaml.Parser())

	// 3) environment variables override (prefix ORDERAPI_, nested with __)
	// e.g. ORDERAPI_MYSQL__DSN, ORDERAPI_REDIS__PASSWORD
	if err := k.Load(env.Provider("ORDERAPI_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "ORDERAPI_")
		s = strings.ReplaceAll(s, "__", ".")
		return strings.ToLower(s)
	}), nil); err != nil {
		return Config{}, fmt.Errorf("env overlay: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	//if c.App.HTTPAddr == "" {
	//	return fmt.Errorf("app.http_addr required")
	//}
	return nil
}
