package serverutils

import (
	"context"
	"errors"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/models/dto"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func ProvideUserCtxMiddleware(ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		combinedCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		requestCtx := c.Context()
		if requestCtx != nil {
			go func() {
				select {
				case <-requestCtx.Done():
					cancel()
					return
				case <-combinedCtx.Done():
					return
				}
			}()
		}

		c.SetUserContext(combinedCtx)
		return c.Next()
	}
}

func RequestLoggingMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		path := c.Path()

		requestUUID := uuid.New()

		requestLogger := log.With().
			Str("method", c.Method()).
			Str("path", path).
			Stringer(string(config.ContextKeyRequestId), requestUUID).
			Logger()

		requestLogger.Debug().
			Int("contentLength", c.Request().Header.ContentLength()).
			Msg("start request")

		ctx = requestLogger.WithContext(ctx)
		ctx = context.WithValue(ctx, config.ContextKeyRequestId, requestUUID)

		defer func(begin time.Time) {
			status := c.Response().StatusCode()
			dur := time.Since(begin)

			requestLogger.Info().
				Dur("executionTime", dur).
				Int("statusCode", status).
				Msgf(
					"[%d %s] %s %s -> took %d ms",
					status,
					http.StatusText(status),
					c.Method(),
					path,
					dur.Milliseconds(),
				)
		}(time.Now())

		c.SetUserContext(ctx)
		return c.Next()
	}
}

func RecoverMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer func() {
			if rec := recover(); rec != nil {
				logRequestErr(c, rec)
			}
		}()

		err := c.Next()

		if err != nil {
			return logRequestErr(c, err)
		}

		return err
	}
}

func logRequestErr(c *fiber.Ctx, anyErr any) error {
	ctx := c.UserContext()

	logMsgBuilder := log.Ctx(ctx).
		Error().
		Stack()

	if runtime.GOOS == "windows" {
		logMsgBuilder.CallerSkipFrame(4)
	} else {
		logMsgBuilder.CallerSkipFrame(2)
	}

	err, ok := anyErr.(error)

	errMsg := ""
	if ok {
		errMsg = err.Error()
		logMsgBuilder = logMsgBuilder.Err(err)
	} else if errStr, ok := anyErr.(string); ok {
		errMsg = errStr
		logMsgBuilder = logMsgBuilder.Interface("err", anyErr)
	}

	logMsgBuilder.Msg("err during request")

	c.Set("Content-Type", "application/json")
	if c.Response().StatusCode() == http.StatusOK {
		c.Status(http.StatusInternalServerError)
	}

	return c.JSON(dto.RequestError{
		Msg:       "err during request",
		Err:       errMsg,
		RequestId: ctx.Value(config.ContextKeyRequestId).(uuid.UUID).String(),
	})
}

// TODO: test
func SecretMiddleware(secret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := strings.Split(c.Get("Authorization"), "Bearer ")
		ctx := c.UserContext()

		if len(authHeader) != 2 {
			log.Ctx(ctx).Debug().Msg("no valid bearer token")
			c.Status(http.StatusUnauthorized)
			return errors.New("no valid bearer token")
		}

		token := authHeader[1]

		if token != secret {
			log.Ctx(ctx).Debug().Msg("no valid token")
			c.Status(http.StatusUnauthorized)
			return errors.New("no valid token")
		}

		return c.Next()
	}
}
