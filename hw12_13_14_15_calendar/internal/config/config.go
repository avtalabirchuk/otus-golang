package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Calendar struct {
	Host        string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port        int    `yaml:"port" env:"PORT" env-default:"8081"`
	GRPCAddress string `yaml:"GRPCAddress" env:"GRPC_ADDRESS"`
	HTTPAddress string `yaml:"HTTPAddress" env:"HTTP_ADDRESS"`
	LogConfig   LogConfig
	DBConfig    DBConfig
}

type Scheduler struct {
	LogConfig   LogConfig
	DBConfig    DBConfig
	QueueConfig QueueConfig
}

type Sender struct {
	LogConfig   LogConfig
	DBConfig    DBConfig
	QueueConfig QueueConfig
}

func readConfig(fpath string, cfg interface{}) error {
	if fpath != "" {
		return cleanenv.ReadConfig(fpath, cfg)
	}
	return cleanenv.ReadEnv(cfg)
}

func NewCalendar(fpath string) (*Calendar, error) {
	var cfg Calendar
	err := readConfig(fpath, &cfg)
	fmt.Printf("ERRRRR %s\n", err)
	fmt.Printf("EXAMPLE_PATH004 %+v \n", cfg)
	return &cfg, err
}

func NewScheduler(fpath string) (*Scheduler, error) {
	cfg := &Scheduler{
		QueueConfig: defaultQueueConfig(),
	}
	return cfg, readConfig(fpath, cfg)
}

func NewSender(fpath string) (*Sender, error) {
	cfg := &Sender{
		QueueConfig: defaultQueueConfig(),
	}
	return cfg, readConfig(fpath, cfg)
}
