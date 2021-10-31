package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
)

var ErrFileLog = errors.New("cannot setup file log")

func getLogLevel(str string) zerolog.Level {
	switch str {
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}

func Init(c *config.Config) (err error) {
	var logInput io.Writer = os.Stderr
	logLevel := getLogLevel(c.LogLevel)
	if c.LogPath != "" {
		f, err := os.OpenFile(c.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			return fmt.Errorf("%s: %w", ErrFileLog, err)
		}
		logInput = zerolog.MultiLevelWriter(f, os.Stderr)
		// defer f.Close()
		// hm... If I close the file, I cannot write there
		// probably, need to close the file, when service is shut down
	}
	log.Logger = zerolog.New(logInput).With().Timestamp().Logger().Level(logLevel)
	return
}
