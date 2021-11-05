package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/scheduler"
)

var cfgPath string

func fatal(err error) {
	log.Fatal().Err(err).Msg("Application cannot start")
}

func init() {
	flag.StringVar(&cfgPath, "config", "", "Calendar Scheduler config")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewScheduler(cfgPath)
	if err != nil {
		fatal(err)
	}
	log.Debug().Msgf("Config inited %+v", cfg)
	if err = logger.Init(&cfg.LogConfig); err != nil {
		fatal(err)
	}
	repo := repository.NewStats(cfg.DBConfig.RepoType, cfg.DBConfig.ItemsPerQuery, cfg.DBConfig.MaxConn)
	if repo == nil {
		fatal(repository.ErrUnSupportedRepoType)
	}
	if err = repo.Connect(ctx, repository.GetSQLDSN(&cfg.DBConfig)); err != nil {
		fatal(err)
	}
	defer repo.Close()

	qCfg := cfg.QueueConfig

	producer := queue.NewProducer(
		qCfg.URI,
		qCfg.QueueName,
		qCfg.ExchangeType,
		qCfg.MaxReconnectAttempts,
		qCfg.ReconnectTimeoutMs,
	)
	app := scheduler.New(repo, producer, qCfg.ScanTimeoutMs)
	if err := app.Run(); err != nil {
		fatal(err)
	}
}
