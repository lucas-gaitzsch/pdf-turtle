package loopback

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/services"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
)

const resourceIdKey = "resourceId"
const bundleIdKey = "bundleId"
const BundlePath = "/bundle"

type Server struct {
	Instance *http.Server
}

func (s *Server) Serve(ctx context.Context) {
	conf := config.Get(ctx)

	servingAddr := fmt.Sprintf("127.0.0.1:%d", conf.LoopbackPort)

	r := mux.NewRouter()

	r.Path(fmt.Sprintf("/preloaded/{%s}", resourceIdKey)).
		Methods(http.MethodGet).
		HandlerFunc(GetPreloadedRessourceHandler).
		Name("Get preloaded resource")

	r.PathPrefix(fmt.Sprintf("%s/{%s}", BundlePath, bundleIdKey)).
		Methods(http.MethodGet).
		HandlerFunc(GetBundleFileHandler).
		Name("Get file of bundle resource")

	s.Instance = &http.Server{
		Handler:      r,
		Addr:         servingAddr,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
		IdleTimeout:  5 * time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go s.listenAndServe()
}

func (s *Server) listenAndServe() {
	log.Debug().Msg("loopback-server: listens on " + s.Instance.Addr)

	if err := s.Instance.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Error().Err(err).Msg("loopback-server serve error")
			panic(err)
		}
	}
}

func (s *Server) Close(ctx context.Context) {
	log.Debug().Msg("loopback-server: shutdown gracefully")

	gracefullyShutdownTimeout := time.Duration(config.Get(ctx).GracefulShutdownTimeoutInSec) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, gracefullyShutdownTimeout)
	defer cancel()
	s.Instance.Shutdown(timeoutCtx)
}

// ### Handler ###

func GetPreloadedRessourceHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// resourceId := vars[resourceIdKey]
	//TODO:
}

func GetBundleFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	bundleIdFromRoute := vars[bundleIdKey]

	urlPathPrefix := BundlePath + "/" + bundleIdFromRoute + "/"

	if r.URL.Path == urlPathPrefix || strings.HasPrefix(urlPathPrefix, r.URL.Path) {
		return
	}

	bundleId, err := uuid.Parse(bundleIdFromRoute)

	if err != nil {
		log.
			Ctx(ctx).
			Error().
			Str("bundleIdFromRoute", bundleIdFromRoute).
			Err(err).
			Msg("cant parse bundle id")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bundleProvider := ctx.Value(config.ContextKeyBundleProviderService).(services.BundleProviderService)

	b, ok := bundleProvider.GetById(bundleId)

	if !ok {
		log.
			Ctx(ctx).
			Error().
			Msg("cant find bundle")
		//TODO:!! http err code
		return
	}

	path := r.URL.Path

	splittedPathByDot := strings.Split(path, ".")
	ext := splittedPathByDot[len(splittedPathByDot)-1]
	mimeType := mime.TypeByExtension(ext)

	fileReader, err := b.GetFileByPath(path[len(urlPathPrefix):])
	if err != nil {
		log.
			Ctx(ctx).
			Error().
			Err(err).
			Msg("cant get file")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fileReader.Close()

	w.Header().Set("Content-type", mimeType)
	_, err = io.Copy(w, fileReader)

	if err != nil {
		log.
			Ctx(ctx).
			Error().
			Err(err).
			Msg("cant respond")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
