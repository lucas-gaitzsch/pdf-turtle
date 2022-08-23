package serverstatic

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/lucas-gaitzsch/pdf-turtle/config"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "github.com/lucas-gaitzsch/pdf-turtle/server/docs"
)

const resourceIdKey = "resourceId"

type Server struct {
	Instance *http.Server
}

func (s *Server) Serve(ctx context.Context) {
	conf := config.Get(ctx)

	servingAddr := fmt.Sprintf("127.0.0.1:%d", conf.Port)

	r := mux.NewRouter()

	r.Path(fmt.Sprintf("/resources/{%s}", resourceIdKey)).
		Methods(http.MethodGet).
		HandlerFunc(GetRessourceHandler).
		Name("Get resource")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet},
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
				log.Error().Err(err).Msg("serverstatic serve error")
				panic(err)
			}
		}
	}()
	log.Info().Msg("serverstatic: listens on " + s.Instance.Addr)
}

func (s *Server) Close(ctx context.Context) {
	log.Info().Msg("serverstatic: shutdown gracefully")
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Get(ctx).GracefulShutdownTimeoutInSec)*time.Second)
	defer cancel()
	s.Instance.Shutdown(timeoutCtx)
}


// ### Handler ###

func GetRessourceHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// resourceId := vars[resourceIdKey]
	//TODO:
}
