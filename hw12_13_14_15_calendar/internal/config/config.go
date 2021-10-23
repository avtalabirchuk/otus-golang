package config

import (
	"context"
	"errors"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Host            string `yaml:"host" config:"required"`
	Port            int    `yaml:"port" config:"required"`
	RepoType        string `yaml:"repoType" config:"required"`
	LogLevel        string `yaml:"logLevel"`
	LogPath         string `yaml:"logPath"`
	DBHost          string `yaml:"dbHost" config:"required"`
	DBPort          int    `yaml:"dbPort" config:"required"`
	DBName          string `yaml:"dbName" config:"required"`
	DBUser          string `yaml:"dbUser" config:"required"`
	DBPass          string `yaml:"dbPass" config:"required"`
	DBMaxConn       int    `yaml:"dBMaxConn"`
	DBItemsPerQuery int    `yaml:"dBItemsPerQuery"`
}

var ErrConfigPath = errors.New("Config path is not provided")

func Read(fpath string) (config *Config, err error) {
	if fpath == "" {
		return nil, ErrConfigPath
	}
	config = &Config{
		Host:            "localhost",
		Port:            8081,
		LogLevel:        "info",
		DBHost:          "localhost",
		DBPort:          5432,
		DBMaxConn:       10,
		DBItemsPerQuery: 100,
	}
	loader := confita.NewLoader(
		file.NewBackend(fpath),
		env.NewBackend(),
	)
	err = loader.Load(context.Background(), config)
	return
}
