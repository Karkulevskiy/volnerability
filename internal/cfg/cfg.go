package cfg

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath string        `json:"storage_path"`
	Address     string        `json:"address"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_timeout"`

	OrchestratorConfig
}

type OrchestratorConfig struct {
	TempDir string `json:"temp_dir"`
}

func MustLoad() *Config {
	cfgPath := "./cfg/cfg.json"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("cfg not found in: %s", cfgPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		l := log.Default()
		cfg = defaultConfig()
		l.Printf("failed read config: %v\n use default: %v", err, cfg)
	}

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		Address:     "0.0.0.0:8082",
		Timeout:     time.Second * 4,
		StoragePath: "./db",
		IdleTimeout: time.Second * 30,
	}
}
