package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucas-gaitzsch/pdf-turtle/models"

	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
)

const TemplateEngineQueryKey = "template-engine"

// RenderPdfFromHtmlHandler godoc
// @Summary      Render PDF from HTML
// @Description  Returns PDF file generated from HTML of body, header and footer
// @Tags         Render HTML
// @Accept       json
// @Produce      application/pdf
// @Param        renderData  body  models.RenderData  true  "Render Data"
// @Success      200         "PDF File"
// @Router       /api/pdf/from/html/render [post]
func RenderPdfFromHtmlHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	data := &models.RenderData{}

	err := c.BodyParser(data)

	if err != nil {
		return err
	}

	pdfService := pdf.NewPdfService(ctx)

	pdfData, err := pdfService.PdfFromHtml(data)

	if err != nil {
		return err
	}

	return writePdf(c, pdfData)
}
