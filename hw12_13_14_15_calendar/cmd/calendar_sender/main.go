package main

import (
	"flag"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/sender"
)

var cfgPath string

func fatal(err error) {
	log.Fatal().Err(err).Msg("Application cannot start")
}

func init() {
	flag.StringVar(&cfgPath, "config", "", "Calendar Sender config")
}

func main() {
	flag.Parse()

	cfg, err := config.NewSender(cfgPath)
	if err != nil {
		fatal(err)
	}
	log.Debug().Msgf("Config inited %+v", cfg)
	if err = logger.Init(&cfg.LogConfig); err != nil {
		fatal(err)
	}

	qCfg := cfg.QueueConfig
	done := make(chan error)

	queueURL := queue.GetRabbitMQURL(qCfg.User, qCfg.Pass, qCfg.Host, strconv.Itoa(qCfg.Port))
	consumer := queue.NewConsumer(
		queueURL,
		qCfg.QueueName,
		qCfg.QueueName,
		qCfg.ExchangeType,
		qCfg.QueueName,
		qCfg.QueueName,
		qCfg.QosPrefetchCount,
		done,
	)
	app := sender.New(consumer, qCfg.ScanTimeout)
	if err := app.Run(done); err != nil {
		fatal(err)
	}
}
