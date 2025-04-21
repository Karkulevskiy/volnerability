package cfg

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Duration time.Duration

type Config struct {
	HttpServer         `yaml:"http_server"`
	OrchestratorConfig `yaml:"orchestrator_config"`
}

type HttpServer struct {
	StoragePath   string     `yaml:"storage_path" env-default:"./internal/db/storage.db"`
	MigrationPath string     `yaml:"migration_path" env-default:"./internal/db/init.sql"`
	Address       string     `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout       Duration   `yaml:"timeout" env-default:"4s"`
	IdleTimeout   Duration   `yaml:"idle_timeout" env-default:"60s"`
	TokenTTL      Duration   `json:"token_ttl" env-default:"1h"`
	GRPCCfg       GRPCConfig `json:"grpc"`
}

type OrchestratorConfig struct {
	TempDir   string `yaml:"temp_dir"`
	TargetDir string `yaml:"target_dir" env-default:"/home"`
	ImageName string `yaml:"image_name" env-default:"code-runner"`
}

type GRPCConfig struct {
	Port    int      `json:"port" env-default:"8085"`
	Timeout Duration `json:"timeout" env-default:"4s"`
}

func MustLoad() *Config {
	const cfgPath = "./internal/cfg/cfg.json"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("cfg not found in: %s", cfgPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return cfg
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

func (d *Duration) SetValue(s string) error {
	// Удаляем кавычки, если они есть
	s = strings.Trim(s, `"`)

	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}
