package cfg

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer         `yaml:"http_server"`
	OrchestratorConfig `yaml:"orchestrator_config"`
}

type HttpServer struct {
	StoragePath string        `yaml:"storage_path" env-default:"./db"`
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type OrchestratorConfig struct {
	TempDir   string `yaml:"temp_dir"`
	TargetDir string `yaml:"target_dir" env-default:"/home"`
	ImageName string `yaml:"image_name" env-default:"code-runner"`
}

func MustLoad() *Config {
	const cfgPath = "./internal/cfg/cfg.yaml"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("cfg not found in: %s", cfgPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return cfg
}
