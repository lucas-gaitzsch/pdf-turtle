package config

import (
	"context"
	"pdf-turtle/utils"

	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevelDebug          bool `arg:"env" default:"true"`
	LogJsonOutput          bool `arg:"env" default:"false"`
	RenderTimeoutInSeconds int  `arg:"env" default:"30"`

	WorkerInstances              int `arg:"env" default:"30"`
	Port                         int `arg:"env" default:"8000"`
	GracefulShutdownTimeoutInSec int `arg:"env" default:"10"`
	MaxBodySizeInMb              int `arg:"env" default:"32"`

	// CachedAssets []string `arg:"env"` //TODO:!
}

func ContextWithConfig(parentCtx context.Context, config Config) context.Context {
	return context.WithValue(parentCtx, ContextKeyConfig, config)
}

func Get(ctx context.Context) Config {
	c, hasConfig := ctx.Value(ContextKeyConfig).(Config)

	if hasConfig {
		return c
	} else {
		log.Warn().Msg("no config was set -> fallback to default")

		c := &Config{}
		utils.ReflectDefaultValues(c)
		return *c
	}

}
