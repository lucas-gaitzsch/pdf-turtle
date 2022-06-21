package renderer

import (
	"context"
	"io"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
)

type HtmlToPdfRendererAbstraction interface {
	RenderHtmlAsPdf(ctx context.Context, data *models.RenderData) (io.Reader, error)
	Close()
}
