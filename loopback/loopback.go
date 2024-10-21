package loopback

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/serverutils"
	"github.com/lucas-gaitzsch/pdf-turtle/services"

	"github.com/rs/zerolog/log"
)

const BundlePath = "/bundle"
const bundleIdKey = "bundleId"

type Server struct {
	Instance *fiber.App
}

func (s *Server) Serve(ctx context.Context) {
	app := fiber.New(fiber.Config{
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
		IdleTimeout:  5 * time.Second,
	})

	app.Use(
		recover.New(),
		serverutils.ProvideUserCtxMiddleware(ctx),
	)

	conf := config.Get(ctx)

	servingAddr := fmt.Sprintf("127.0.0.1:%d", conf.LoopbackPort)

	app.
		Get(fmt.Sprintf("%s/:%s/+", BundlePath, bundleIdKey), GetBundleFileHandler).
		Name("Get file of bundle resource")

	s.Instance = app

	go s.listenAndServe(servingAddr)
}

func (s *Server) listenAndServe(servingAddr string) {
	log.Debug().Msg("loopback-server: listens on " + servingAddr)

	if err := s.Instance.Listen(servingAddr); err != nil {
		if err != http.ErrServerClosed {
			log.Error().Err(err).Msg("loopback-server serve error")
			panic(err)
		}
	}
}

func (s *Server) Close(ctx context.Context) {
	gracefullyShutdownTimeout := time.Duration(config.Get(ctx).GracefulShutdownTimeoutInSec) * time.Second
	s.Instance.ShutdownWithTimeout(gracefullyShutdownTimeout)
}

// ### Handler ###

func GetBundleFileHandler(c *fiber.Ctx) error {
	if c.Method() != http.MethodGet {
		return c.SendStatus(http.StatusMethodNotAllowed)
	}

	ctx := c.UserContext()

	bundleIdFromRoute := c.Params(bundleIdKey)

	urlPathPrefix := BundlePath + "/" + bundleIdFromRoute + "/"

	path := c.Path()

	if path == urlPathPrefix || strings.HasPrefix(urlPathPrefix, path) {
		return c.SendStatus(http.StatusInternalServerError)
	}

	bundleId, err := uuid.Parse(bundleIdFromRoute)

	if err != nil {
		log.
			Ctx(ctx).
			Error().
			Str("bundleIdFromRoute", bundleIdFromRoute).
			Err(err).
			Msg("cant parse bundle id")

		return c.SendStatus(http.StatusInternalServerError)
	}

	bundleProvider := ctx.Value(config.ContextKeyBundleProviderService).(services.BundleProviderService)

	b, ok := bundleProvider.GetById(bundleId)

	if !ok {
		log.
			Ctx(ctx).
			Error().
			Msg("cant find bundle")

		return c.SendStatus(http.StatusNotFound)
	}

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

		return c.SendStatus(http.StatusInternalServerError)
	}
	defer fileReader.Close()

	c.Set(fiber.HeaderContentType, mimeType)
	return c.SendStream(fileReader)
}
