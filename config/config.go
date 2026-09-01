package config

import (
	"flag"
	"os"

	grpcsrv "github.com/Cheasezz/fileService/internal/grpc"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string         `yaml:"env" env-default:"local"`
	GRPC grpcsrv.Config `yaml:"grpc" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

// fetch config from command line flag, or from os env
// Priority: flag > env > default
func fetchConfigPath() string {
	cfgPath := ""

	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = os.Getenv("FILE_SERVICE_CONFIG_PATH")
	}

	return cfgPath
}
