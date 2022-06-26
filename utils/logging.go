package utils

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ContextKey string

const contextKeySkipFrames = ContextKey("contextKeySkipFrames")

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

	skipFrames, ok := ctx.Value(contextKeySkipFrames).(int)
	if !ok {
		skipFrames = 1
	}

	logger.
		Debug().
		Dur("executionTime", duration).
		CallerSkipFrame(skipFrames).
		Msgf("%s: %d ms", msg, duration.Milliseconds())
}

func LogExecutionTimeWithResult[R any, E error](msg string, ctx context.Context, f func() (R, E)) (R, E) {
	var res R
	var err E

	LogExecutionTime(msg, context.WithValue(ctx, contextKeySkipFrames, 2), func() {
		res, err = f()
	})

	return res, err
}
