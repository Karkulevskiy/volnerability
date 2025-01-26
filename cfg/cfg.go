package cfg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	StoragePath string `json:"storage_path"`
	HTTPServer
}

type HTTPServer struct {
	Address     string        `json:"address"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	User        string        `json:"user"`
	Password    string        `json:"password"`
}

func MustLoad() *Config {
	viper.SetConfigName("cfg.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed read config: %v", err)
	}

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("failed parse config: %v", err)
	}

	return cfg
}
