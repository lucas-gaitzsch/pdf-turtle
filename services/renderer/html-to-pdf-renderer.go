package renderer

import (
	"context"
	"io"
	"pdf-turtle/models"
)

type HtmlToPdfRendererAbstraction interface {
	RenderHtmlAsPdf(ctx context.Context, data *models.RenderData) (io.Reader, error)
	Close()
}
