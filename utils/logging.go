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

	skipFrames := 1

	if ctx != nil {
		if sf, ok := ctx.Value(contextKeySkipFrames).(int); ok {
			skipFrames = sf
		}
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

	var innerCtx context.Context
	if ctx != nil {
		innerCtx = context.WithValue(ctx, contextKeySkipFrames, 2)
	}

	LogExecutionTime(msg, innerCtx, func() {
		res, err = f()
	})

	return res, err
}
