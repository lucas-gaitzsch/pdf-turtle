package config

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lucas-gaitzsch/pdf-turtle/utils"

	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevelDebug          bool `arg:"--logDebug,env:LOG_LEVEL_DEBUG" default:"false" help:"Debug log level active"`
	LogJsonOutput          bool `arg:"--logJsonOutput,env:LOG_JSON_OUTPUT" default:"false" help:"Json log output"`
	RenderTimeoutInSeconds int  `arg:"--renderTimeout,env:RENDER_TIMEOUT" default:"30" help:"Render timeout in seconds"`
	WorkerInstances        int  `arg:"--workerInstances,env:WORKER_INSTANCES" default:"30"`

	Port                         int    `arg:"env" default:"8000" help:"Server port"`
	GracefulShutdownTimeoutInSec int    `arg:"--GracefulShutdownTimeout,env:GRACEFUL_SHUTDOWN_TIMEOUT" default:"10" help:"Graceful server shutdown timeout in seconds"`
	MaxBodySizeInMb              int    `arg:"--maxBodySize,env:MAX_BODY_SIZE" default:"32" help:"Max body size in megabyte"`
	ServePlayground              bool   `arg:"--servePlayground,env:SERVE_PLAYGROUND" default:"false" help:"Serve playground from path './static-files/playground/'"`
	Secret                       string `arg:"env" default:"" help:"Secret used as bearer token"`
	NoSandbox                    bool   `arg:"--no-sandbox,env:NO_SANDBOX" default:"false" help:"Disable chromium sandbox"`

	PreloadedAssets []string `arg:"env" help:"Preload assets on startup. Example:'bar.js:https://foo.com/bar.js'"` //TODO:!

	LoopbackPort int `arg:"env" default:"8001" help:"Loopback-Server port"`

	EnableUrlRender bool   `arg:"--enableUrlRender,env:ENABLE_URL_RENDER" default:"false" help:"Enable URL render endpoint for fetching and rendering HTML from URLs"`
	ProxyUrl        string `arg:"--proxyUrl,env:PROXY_URL" default:"" help:"HTTP proxy URL for outbound requests (e.g., http://proxy.example.com:8080)"`
	proxyUrlParsed  *url.URL
}

// Validate validates the configuration and returns an error if invalid
func (c *Config) Validate() error {
	if c.ProxyUrl != "" {
		proxyURL, err := url.Parse(c.ProxyUrl)
		if err != nil {
			return fmt.Errorf("invalid proxy url: %w", err)
		}
		c.proxyUrlParsed = proxyURL
	}
	return nil
}

// GetProxyUrl returns the parsed proxy URL if configured
func (c *Config) GetProxyUrl() *url.URL {
	return c.proxyUrlParsed
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
