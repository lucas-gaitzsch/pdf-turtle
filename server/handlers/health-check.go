package handlers

import (
	"net/http"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
)

// Health Endpoint (liveness) godoc
// @Summary      Liveness probe for this service
// @Tags         Internals
// @Accept       multipart/form-data
// @Produce      text/plain
// @Success      200
// @Router       /api/health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	testHtml := "health"

	data := &models.RenderData{
		Html: &testHtml,
	}

	pdfService := pdf.NewPdfService(ctx)

	_, err := pdfService.PdfFromHtml(data)

	if err != nil {
		panic(err)
	}
}
