package utils

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogExecutionTime(msg string, ctx context.Context, f func()) {
	start := time.Now()

	f()

	duration := time.Since(start)

	var logger *zerolog.Logger
	if ctx == nil {
		logger = &log.Logger
	} else {
		logger = log.Ctx(ctx)
	}

	logger.
		Debug().
		Dur("executionTime", duration).
		CallerSkipFrame(1).
		Msgf("%s: %d ms", msg, duration.Milliseconds())
}

func LogExecutionTimeWithResult[R any, E error](msg string, ctx context.Context, f func() (R, E)) (R, E) {
	var res R
	var err E

	LogExecutionTime(msg, ctx, func() {
		res, err = f()
	})

	return res, err
}
