package main

import (
	"context"
	"os"
	"os/signal"
	"pdf-turtle/config"
	"pdf-turtle/server"
	"pdf-turtle/services/assetsprovider"
	"pdf-turtle/services/renderer"
	"pdf-turtle/utils/logging"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/alexflint/go-arg"
)

func initConfigByArgs(ctx context.Context) context.Context {
	var c config.Config
	arg.MustParse(&c)

	return config.ContextWithConfig(ctx, c)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = initConfigByArgs(ctx)

	logging.InitLogger(ctx)

	log.Info().Msg("Hey dude 👋 .. I am Karl, your turtle for today 🐢")

	// init services
	servicesCtx := ctx

	rendererService := renderer.NewRendererBackgroundService(ctx)
	servicesCtx = context.WithValue(servicesCtx, config.ContextKeyRendererService, rendererService)

	assetsProviderService := assetsprovider.NewAssetsProviderService()
	servicesCtx = context.WithValue(servicesCtx, config.ContextKeyAssetsProviderService, assetsProviderService)

	// init server
	srv := server.Server{}
	srv.Serve(servicesCtx)

	// listen to os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	<-osSignals

	// cleanup resources
	log.Info().Msg("shutting service down...")
	srv.Close(ctx)
	rendererService.Close()
	cancel()

	// exit
	log.Info().Msg("🐢 bye bye dude.")
	os.Exit(0)
}
