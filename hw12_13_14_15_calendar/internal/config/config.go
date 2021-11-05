package config

import (
	"context"
	"errors"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type Calendar struct {
	Host        string    `yaml:"host" config:"required"`
	Port        int       `yaml:"port" config:"required"`
	GRPCAddress string    `yaml:"grpcAddress"`
	HTTPAddress string    `yaml:"httpAddress"`
	LogConfig   LogConfig `yaml:"logConfig"`
	DBConfig    DBConfig  `yaml:"dbConfig"`
}

type Scheduler struct {
	LogConfig   LogConfig   `yaml:"logConfig"`
	DBConfig    DBConfig    `yaml:"dbConfig"`
	QueueConfig QueueConfig `yaml:"queueConfig"`
}

type Sender struct {
	LogConfig   LogConfig   `yaml:"logConfig"`
	DBConfig    DBConfig    `yaml:"dbConfig"`
	QueueConfig QueueConfig `yaml:"queueConfig"`
}

var ErrConfigPath = errors.New("config path is not provided")

func NewCalendar(fpath string) (config *Calendar, err error) {
	if fpath == "" {
		return nil, ErrConfigPath
	}
	config = &Calendar{
		Host:      "localhost",
		Port:      8081,
		LogConfig: defaultLogConfig(),
		DBConfig:  defaultDBConfig(),
	}
	loader := confita.NewLoader(
		file.NewBackend(fpath),
		env.NewBackend(),
	)
	err = loader.Load(context.Background(), config)
	return
}

func NewScheduler(fpath string) (config *Scheduler, err error) {
	if fpath == "" {
		return nil, ErrConfigPath
	}
	config = &Scheduler{
		LogConfig:   defaultLogConfig(),
		DBConfig:    defaultDBConfig(),
		QueueConfig: defaultQueueConfig(),
	}
	loader := confita.NewLoader(
		file.NewBackend(fpath),
		env.NewBackend(),
	)
	err = loader.Load(context.Background(), config)
	return
}

func NewSender(fpath string) (config *Sender, err error) {
	if fpath == "" {
		return nil, ErrConfigPath
	}
	config = &Sender{
		LogConfig:   defaultLogConfig(),
		QueueConfig: defaultQueueConfig(),
	}
	loader := confita.NewLoader(
		file.NewBackend(fpath),
		env.NewBackend(),
	)
	err = loader.Load(context.Background(), config)
	return
}
