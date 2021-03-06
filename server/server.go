package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/server/handlers"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "github.com/lucas-gaitzsch/pdf-turtle/server/docs"

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
// @description     A painless HTML to PDF rendering service. Generate PDF reports and documents from HTML templates or raw HTML.
// @contact.name    Lucas Gaitzsch
// @contact.email   lucas@gaitzsch.dev
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes         http
// @basePath        /api
func (s *Server) Serve(ctx context.Context) {
	conf := config.Get(ctx)

	servingAddr := fmt.Sprintf(":%d", conf.Port)
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

	api.Use(
		maxBodySizeMiddleware(conf.MaxBodySizeInMb),
		loggingMiddleware(),
		recoverMiddleware(),
	)

	if conf.Secret != "" {
		api.Use(secretMiddleware(conf.Secret))
	}

	// Swagger
	r.PathPrefix(swaggerRoute).Handler(httpSwagger.WrapHandler)

	log.
		Info().
		Str("url", fmt.Sprintf("%s%s/index.html", localUrl, swaggerRoute)).
		Msg("serving open-api (swagger) description")

	if conf.ServePlayground {
		log.
			Info().
			Str("url", localUrl).
			Msg("serving playground")

		r.PathPrefix("/assets").Handler(http.FileServer(http.Dir(config.PathStaticExternPlayground)))

		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/favicon.ico" {
				http.ServeFile(w, r, config.PathStaticExternPlayground+"favicon.ico")
			} else {
				http.ServeFile(w, r, config.PathStaticExternPlayground+"index.html")
			}
		})
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	})

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
