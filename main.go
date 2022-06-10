package main

import (
	"context"
	"os"
	"os/signal"
	"pdf-turtle/config"
	"pdf-turtle/server"
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

	pdfService := renderer.NewRendererBackgroundService(ctx)

	// init server
	serverCtx := context.WithValue(ctx, config.ContextKeyPdfService, pdfService)
	srv := server.Server{}
	srv.Serve(serverCtx)

	// listen to os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	<-osSignals

	// cleanup resources
	log.Info().Msg("shutting service down...")
	srv.Close(ctx)
	pdfService.Close()
	cancel()

	// exit
	log.Info().Msg("🐢 bye bye dude.")
	os.Exit(0)
}
