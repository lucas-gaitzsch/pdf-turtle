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

func LogExecutionTimeWithResult[R1 any](msg string, ctx context.Context, f func() R1) R1 {
	var res1 R1

	var innerCtx context.Context
	if ctx != nil {
		innerCtx = context.WithValue(ctx, contextKeySkipFrames, 2)
	}

	LogExecutionTime(msg, innerCtx, func() {
		res1 = f()
	})

	return res1
}

func LogExecutionTimeWithResults[R1 any, R2 any](msg string, ctx context.Context, f func() (R1, R2)) (R1, R2) {
	var res1 R1
	var res2 R2

	var innerCtx context.Context
	if ctx != nil {
		innerCtx = context.WithValue(ctx, contextKeySkipFrames, 2)
	}

	LogExecutionTime(msg, innerCtx, func() {
		res1, res2 = f()
	})

	return res1, res2
}
