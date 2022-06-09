package main

import (
	"context"
	"os"
	"os/signal"
	"pdf-turtle/config"
	"pdf-turtle/server"
	"pdf-turtle/services/renderer"
	"pdf-turtle/utils/logging"

	"github.com/rs/zerolog/log"

	"github.com/alexflint/go-arg"
)

func initConfigByArgs(ctx context.Context) context.Context {
	var c config.Config
	arg.MustParse(&c)

	return config.ContextWithConfig(ctx, c)
}

func main() {
	ctx := initConfigByArgs(context.Background())
	
	logging.InitLogger(ctx)

	log.Info().Msg("Hey dude 👋 .. I am Karl, your turtle for today 🐢")

	ctx, cancel := context.WithCancel(ctx)

	pdfService := renderer.NewRendererBackgroundService(ctx)

	serverCtx := context.WithValue(ctx, config.ContextKeyPdfService, pdfService)

	srv := server.Server{}
	srv.Serve(serverCtx)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	
	<-sigint
	
	log.Info().Msg("shutting service down...")
	srv.Close(ctx)
	pdfService.Close()
	cancel()

	log.Info().Msg("🐢 bye bye dude.")

	os.Exit(0)
}
