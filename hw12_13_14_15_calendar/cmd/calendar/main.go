package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/calendar"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

var ErrUnSupportedRepoType = errors.New("unsupported repository type")

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

	cfg, err := config.Read(cfgPath)
	if err != nil {
		fatal(err)
	}
	log.Debug().Msgf("Config inited %+v", cfg)
	err = logger.Init(cfg)
	if err != nil {
		fatal(err)
	}
	repo := repository.New(cfg.RepoType)
	if repo == nil {
		fatal(ErrUnSupportedRepoType)
	}
	err = repo.Connect(ctx, cfg)
	if err != nil {
		fatal(err)
	}
	defer repo.Close()

	app, err := calendar.New(repo)
	if err != nil {
		fatal(err)
	}
	err = app.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		fatal(err)
	}
}
