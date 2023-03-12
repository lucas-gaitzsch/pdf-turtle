package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/server/handlers"
	"github.com/lucas-gaitzsch/pdf-turtle/serverutils"

	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	fiberSwagger "github.com/swaggo/fiber-swagger" // fiber-swagger middleware

	_ "github.com/lucas-gaitzsch/pdf-turtle/server/docs"
)

const (
	swaggerRoute = "/swagger"
)

type Server struct {
	Instance *fiber.App
}

// @title          PdfTurtle API
// @version        1.1
// @description    A painless HTML to PDF rendering service. Generate PDF reports and documents from HTML templates or raw HTML.
// @contact.name   Lucas Gaitzsch
// @contact.email  lucas@gaitzsch.dev
// @license.name   AGPL-3.0
// @license.url    https://github.com/lucas-gaitzsch/pdf-turtle/blob/main/LICENSE

func (s *Server) Serve(ctx context.Context) {
	conf := config.Get(ctx)

	app := fiber.New(fiber.Config{
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
		IdleTimeout:  60 * time.Second,
		BodyLimit: conf.MaxBodySizeInMb * 1024 * 1024,
    })

	app.Use(
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders:  "*",
			AllowMethods: strings.Join([]string{http.MethodGet, http.MethodPost}, ","),
			AllowCredentials: true,
		}),
		recover.New(),
		serverutils.ProvideUserCtxMiddleware(ctx),
	)

	servingAddr := fmt.Sprintf(":%d", conf.Port)
	localUrl := fmt.Sprintf("http://localhost%s", servingAddr)

	app.
		Get("/health", handlers.HealthCheckHandler).
		Name("Liveness probe")
	app.
		Get("/api/health", handlers.HealthCheckHandler).
		Name("Liveness probe")

	api := app.Group("/api")

	api.Use(
		serverutils.RequestLoggingMiddleware(),
		serverutils.RecoverMiddleware(),
	)

	api.Post("/pdf/from/html/render", handlers.RenderPdfFromHtmlHandler).
		Name("Render PDF from HTML")

	api.Post("/pdf/from/html-template/render", handlers.RenderPdfFromHtmlFromTemplateHandler).
		Name("Render PDF from HTML template")

	api.Post("/pdf/from/html-template/test", handlers.TestHtmlTemplateHandler).
		Name("Test HTML template")

	api.Post("/pdf/from/html-bundle/render", handlers.RenderBundleHandler).
		Name("Render PDF from HTML-Bundle")

	if conf.Secret != "" {
		api.Use(serverutils.SecretMiddleware(conf.Secret))
	}

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.
		Info().
		Str("url", fmt.Sprintf("%s%s/index.html", localUrl, swaggerRoute)).
		Msg("serving open-api (swagger) description")

	// Serve playground vue frontend if required
	if conf.ServePlayground {
		log.
			Info().
			Str("url", localUrl).
			Msg("serving playground")
		servePlaygroundFronted(app)
	}

	s.Instance = app

	go s.listenAndServe(servingAddr)
}

func (s *Server) listenAndServe(servingAddr string) {
	log.Info().Msg("server: listens on " + servingAddr)

	if err := s.Instance.Listen(servingAddr); err != nil {
		if err != http.ErrServerClosed {
			log.Error().Err(err).Msg("loopback-server serve error")
			panic(err)
		}
	}
}

func (s *Server) Close(ctx context.Context) {
	log.Info().Msg("server: shutdown gracefully")
	gracefullyShutdownTimeout := time.Duration(config.Get(ctx).GracefulShutdownTimeoutInSec) * time.Second
	s.Instance.ShutdownWithTimeout(gracefullyShutdownTimeout)
}

func servePlaygroundFronted(app *fiber.App) {
	app.Static("/assets", config.PathStaticExternPlayground)
	app.Static("/favicon.ico", config.PathStaticExternPlayground+"favicon.ico")
	app.Static("*", config.PathStaticExternPlayground+"index.html")
}
