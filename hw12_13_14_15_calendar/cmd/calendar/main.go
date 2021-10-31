package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

var cfgPath string

func fatal(err error) {
	log.Fatal().Err(err).Msg("Application cannot start")
}

func init() {
	flag.StringVar(&cfgPath, "config", "", "Calendar config")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewCalendar(cfgPath)
	if err != nil {
		fatal(err)
	}
	log.Debug().Msgf("Config inited %+v", cfg)

	if err := logger.Init(&cfg.LogConfig); err != nil {
		fatal(err)
	}

	repo := repository.NewCRUD(cfg.DBConfig.RepoType, cfg.DBConfig.ItemsPerQuery, cfg.DBConfig.MaxConn)
	if repo == nil {
		fatal(repository.ErrUnSupportedRepoType)
	}

	if err = repo.Connect(ctx, repository.GetSQLDSN(&cfg.DBConfig)); err != nil {
		fatal(err)
	}
	defer repo.Close()

	app, err := app.New(cfg, repo)
	if err != nil {
		fatal(err)
	}

	errCh := make(chan error)
	doneCh := make(chan bool)
	go app.Run(errCh, doneCh)

	for {
		select {
		case <-doneCh:
			return
		case err := <-errCh:
			log.Error().Msgf("%s", err)
		}
	}
}
