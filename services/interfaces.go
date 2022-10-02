package services

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services/bundles"
)

type AssetsProviderService interface {
	GetMergedCss() *string
	GetCssByKey(key string) (css *string, ok bool)
}

type BundleProviderService interface {
	Provide(bundle *bundles.Bundle) (id uuid.UUID, cleanup bundles.CleanupFunc)
	Remove(id uuid.UUID)
	GetById(id uuid.UUID) (bundles.BundleReader, bool)
}

type RendererBackgroundService interface {
	Init(outerCtx context.Context)
	RenderAndReceive(job models.Job) (io.Reader, error)
	Close()
}
