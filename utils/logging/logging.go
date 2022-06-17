package logging

import (
	"context"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"pdf-turtle/config"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(ctx context.Context) {
	var w io.Writer

	conf := config.Get(ctx)

	if conf.LogJsonOutput {
		w = os.Stdout
	} else {
		w = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	loggerContext := zerolog.
		New(w).
		With().
		Timestamp()

	minLevelDebug := conf.LogLevelDebug

	if minLevelDebug {
		loggerContext = loggerContext.Caller()
	}

	logger := loggerContext.Logger()

	if !minLevelDebug {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	log.Logger = logger
}

func InitTestLogger(t *testing.T) {
	loggerContext := zerolog.New(zerolog.ConsoleWriter{
		Out: zerolog.TestWriter{
			T: t,
		},
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Timestamp()

	loggerContext = loggerContext.Caller()

	logger := loggerContext.Logger()

	logger.Level(zerolog.InfoLevel)

	log.Logger = logger
}

func SetNullLogger() {
	log.Logger = zerolog.
		New(ioutil.Discard).
		With().
		Logger().
		Level(zerolog.Disabled)
}
