package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/models/dto"

	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
	"github.com/lucas-gaitzsch/pdf-turtle/services/templating/templateengines"
)

// RenderPdfFromHtmlFromTemplateHandler godoc
// @Summary      Render PDF from HTML template
// @Description  Returns PDF file generated from HTML template plus model of body, header and footer
// @Tags         Render HTML-Template
// @Accept       json
// @Produce      application/pdf
// @Param        renderTemplateData  body      models.RenderTemplateData  true  "Render Data"
// @Success      200                 "PDF File"
// @Router       /api/pdf/from/html-template/render [post]
func RenderPdfFromHtmlFromTemplateHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	templateData := &models.RenderTemplateData{}

	err := c.BodyParser(templateData)

	if err != nil {
		return err
	}

	pdfService := pdf.NewPdfService(ctx)

	pdfData, err := pdfService.PdfFromHtmlTemplate(templateData)

	if err != nil {
		return err
	}

	return writePdf(c, pdfData)
}

// TestHtmlTemplateHandler godoc
// @Summary      Test HTML template matching model
// @Description  Returns information about matching model data to template
// @Tags         Render HTML-Template
// @Accept       json
// @Produce      json
// @Param        renderTemplateData  body  models.RenderTemplateData  true  "Render Data"
// @Success      200                 {object}  dto.TemplateTestResult
// @Router       /api/pdf/from/html-template/test [post]
func TestHtmlTemplateHandler(c *fiber.Ctx) error {
	response := dto.TemplateTestResult{}

	templateData := &models.RenderTemplateData{}

	err := c.BodyParser(templateData)

	if err == nil {

		templateData.ParseJsonModelDataFromDoubleEncodedString()

		templateEngine, found := templateengines.GetTemplateEngineByKey(templateData.TemplateEngine)

		templateengines.LogParsedTemplateEngine(templateData.TemplateEngine, templateEngine, found)

		bodyErr := templateEngine.Test(templateData.HtmlTemplate, templateData.Model)
		if bodyErr != nil {
			strErr := bodyErr.Error()
			response.BodyTemplateError = &strErr
		}
	}

	if response.BodyTemplateError == nil && response.HeaderTemplateError == nil && response.FooterTemplateError == nil {
		response.IsValid = true
	} else {
		response.IsValid = false
	}

	return c.JSON(response)
}
