package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/models/dto"

	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
	"github.com/lucas-gaitzsch/pdf-turtle/services/templating/templateengines"
)

// RenderPdfFromHtmlFromTemplateHandler godoc
// @Summary      Render PDF from HTML template
// @Description  Returns PDF file generated from HTML template plus model of body, header and footer
// @Tags         render html-template
// @Accept       json
// @Produce      application/pdf
// @Param        renderTemplateData  body      models.RenderTemplateData  true  "Render Data"
// @Success      200                 "PDF File"
// @Router       /pdf/from/html-template/render [post]
func RenderPdfFromHtmlFromTemplateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	templateData := &models.RenderTemplateData{}

	err := json.NewDecoder(r.Body).Decode(templateData)

	if err != nil {
		panic(err)
	}

	pdfService := pdf.NewPdfService(ctx)

	pdfData, err := pdfService.PdfFromHtmlTemplate(templateData)

	if err != nil {
		panic(err)
	}

	if err := writePdf(ctx, w, pdfData); err != nil {
		panic(err)
	}
}

// TestHtmlTemplateHandler godoc
// @Summary      Test HTML template matching model
// @Description  Returns information about matching model data to template
// @Tags         test html-template
// @Accept       json
// @Produce      json
// @Param        renderTemplateData  body  models.RenderTemplateData  true  "Render Data"
// @Success      200                 {object}  dto.TemplateTestResult
// @Router       /pdf/from/html-template/test [post]
func TestHtmlTemplateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response := dto.TemplateTestResult{}

	templateData := &models.RenderTemplateData{}

	err := json.NewDecoder(r.Body).Decode(templateData)

	if err == nil {

		templateData.ParseJsonModelDataFromDoubleEncodedString()

		templateEngine := templateengines.GetTemplateEngineByKey(templateData.TemplateEngine)

		bodyErr := templateEngine.Test(templateData.HtmlTemplate, templateData.Model)
		if bodyErr != nil {
			strErr := bodyErr.Error()
			response.BodyTemplateError = &strErr
		}

		headerErr := templateEngine.Test(&templateData.HeaderHtmlTemplate, templateData.HeaderModel)
		if headerErr != nil {
			strErr := headerErr.Error()
			response.HeaderTemplateError = &strErr
		}

		footerErr := templateEngine.Test(&templateData.FooterHtmlTemplate, templateData.FooterModel)
		if footerErr != nil {
			strErr := footerErr.Error()
			response.FooterTemplateError = &strErr
		}
	}

	if response.BodyTemplateError == nil && response.HeaderTemplateError == nil && response.FooterTemplateError == nil {
		response.IsValid = true
	} else {
		response.IsValid = false
	}

	if err := writeJson(ctx, w, response); err != nil {
		panic(err)
	}
}
