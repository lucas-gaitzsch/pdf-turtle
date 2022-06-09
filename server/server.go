package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"pdf-turtle/config"
	"pdf-turtle/server/handlers"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "pdf-turtle/server/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	swaggerRoute = "/swagger"
)

type Server struct {
	Instance *http.Server
}

// @title           PdfTurtle API
// @version         1.0
// @description     #TODO
// @termsOfService  #TODO
// @contact.name    Lucas Gaitzsch
// @contact.email   lucas@gaitzsch.dev
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes         http
// @basePath        /api
func (s *Server) Serve(ctx context.Context) {
	servingAddr := fmt.Sprintf(":%d", config.Get(ctx).Port)
	localUrl := fmt.Sprintf("http://localhost%s", servingAddr)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.Path("/pdf/from/html/render").
		Methods(http.MethodPost).
		HandlerFunc(handlers.RenderPdfFromHtmlHandler).
		Name("Render PDF from HTML")

	api.Path("/pdf/from/html-template/render").
		Methods(http.MethodPost).
		HandlerFunc(handlers.RenderPdfFromHtmlFromTemplateHandler).
		Name("Render PDF from HTML template")

	api.Path("/pdf/from/html-template/test").
		Methods(http.MethodPost).
		HandlerFunc(handlers.TestHtmlTemplateHandler).
		Name("Test HTML template")

	// Swagger
	r.PathPrefix(swaggerRoute).Handler(httpSwagger.WrapHandler)

	log.
		Info().
		Str("url", fmt.Sprintf("%s%s/index.html", localUrl, swaggerRoute)).
		Msg("serving open-api (swagger) description")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	})

	r.Use(
		maxBodySizeMiddleware(config.Get(ctx).MaxBodySizeInMb),
		loggingMiddleware(),
		recoverMiddleware(),
	)

	handler := c.Handler(r)

	s.Instance = &http.Server{
		Handler:      handler,
		Addr:         servingAddr,
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
		IdleTimeout:  60 * time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err := s.Instance.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Error().Err(err).Msg("server serve error")
				panic(err)
			}
		}
	}()
	log.Info().Msg("server: listens on " + s.Instance.Addr)
}

func (s *Server) Close(ctx context.Context) {
	log.Info().Msg("server: shutdown gracefully")
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Get(ctx).GracefulShutdownTimeoutInSec)*time.Second)
	defer cancel()
	s.Instance.Shutdown(timeoutCtx)
}
