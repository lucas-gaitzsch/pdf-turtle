package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
	"github.com/rs/zerolog/log"
)

// Health Endpoint (liveness) godoc
// @Summary      Liveness probe for this service
// @Tags         Internals
// @Accept       multipart/form-data
// @Produce      text/plain
// @Success      200
// @Router       /health [get]
func HealthCheckHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	log.Ctx(ctx).Debug().Msg("execute health check / liveness probe")

	testHtml := "health"

	data := &models.RenderData{
		Html: &testHtml,
	}

	pdfService := pdf.NewPdfService(ctx)

	_, err := pdfService.PdfFromHtml(data)

	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}
