package main

import (
	"flag"
	"fmt"
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

	consumer := queue.NewConsumer(
		fmt.Sprintf("amqp://%s:%s@%s:%s/", qCfg.User, qCfg.Pass, qCfg.Host, strconv.Itoa(qCfg.Port)),
		qCfg.QueueName,
		qCfg.ExchangeType,
		qCfg.QosPrefetchCount,
		qCfg.MaxReconnectAttempts,
		qCfg.ReconnectTimeoutMs,
	)
	app := sender.New(consumer, qCfg.ScanTimeoutMs)
	if err := app.Run(); err != nil {
		fatal(err)
	}
}
