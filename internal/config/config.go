package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpSever struct {
	Address string
}

// Used cleanenv package in go to mapped env variables or yaml variables with structs
type Config struct {
	Env         string `yaml:"env" env:"env"`
	StoragePath string `yaml:"storage_path" env:"storage_path"`
	HttpSever   `yaml:"http_server" env:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	// Get all env files to get local.yaml file
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// Check if the files is pass as flag from command
		flags := flag.String("config", "", "path to the configuration file")

		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist %s", configPath)
	}

	var cfg Config

	// Use ReadConfig method from cleanenv library and pass pointer to store config path
	err1 := cleanenv.ReadConfig(configPath, &cfg)
	if err1 != nil {
		log.Fatalf("can not read config file %s", err1.Error())
	}

	return &cfg
}
