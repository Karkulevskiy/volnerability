package cfg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	StoragePath string        `json:"storage_path"`
	Address     string        `json:"address"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	TokenTTL 	time.Duration `json:"token_ttl"`
	GRPCCfg 	GRPCConfig	  `json:"grpc"`
}

type GRPCConfig struct{
	Port int 			  `json:"port"`
	Timeout time.Duration `json:"timeout"`
}

func MustLoad() *Config {
	viper.SetConfigFile("./internal/cfg/cfg.json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("cfg file was not found, using default settings")
			return defaultConfig()
		} else {
			log.Fatalf("failed read config: %v", err)
		}
	}

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("failed parse config: %v", err)
	}

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		Address:     "0.0.0.0:8082",
		Timeout:     time.Second * 4,
		StoragePath: "./db",
		IdleTimeout: time.Second * 30,
		TokenTTL: time.Hour * 1,
		GRPCCfg: GRPCConfig{
			Port:8085,
			Timeout: time.Second * 5,
		},
	}
}
