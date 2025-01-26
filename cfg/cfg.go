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
}

func MustLoad() *Config {
	viper.SetConfigName("cfg")
	viper.SetConfigType("json")
	viper.AddConfigPath("cfg\\")

	viper.SetDefault("address", "0.0.0.0:8082") // TODO добавить дефолтные настройки

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed read config: %v", err)
	}

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("failed parse config: %v", err)
	}

	return cfg
}
