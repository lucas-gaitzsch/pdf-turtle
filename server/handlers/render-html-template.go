package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"pdf-turtle/models"
	"pdf-turtle/models/dto"
	"pdf-turtle/templating"
	"pdf-turtle/utils"
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
	pdfService := getPdfService(ctx)

	templateData := &models.RenderTemplateData{}

	err := json.NewDecoder(r.Body).Decode(templateData)

	if err != nil {
		panic(err)
	}

	templateData.ParseJsonModelDataFromDoubleEncodedString()

	templateEngine := templating.GetTemplateEngineByKey(templateData.TemplateEngine)

	data := &models.RenderData{
		RenderOptions: templateData.RenderOptions,
	}

	utils.LogExecutionTime("exec template", &ctx, func() {
		data.BodyHtml, err = templateEngine.Execute(templateData.BodyHtmlTemplate, templateData.BodyModel)
		if err != nil {
			panic(err)
		}

		headerHtml, err := templateEngine.Execute(&templateData.HeaderHtmlTemplate, templateData.HeaderModel)
		if err != nil {
			panic(err)
		}

		footerHtml, err := templateEngine.Execute(&templateData.FooterHtmlTemplate, templateData.FooterModel)
		if err != nil {
			panic(err)
		}

		data.HeaderHtml = *headerHtml
		data.FooterHtml = *footerHtml
	})

	pdfData, err := utils.LogExecutionTimeWithResult("render pdf", &ctx, func() (io.Reader, error) {
		return pdfService.RenderAndReceive(*models.NewJob(ctx, data))
	})

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

		templateEngine := templating.GetTemplateEngineByKey(templateData.TemplateEngine)

		bodyErr := templateEngine.Test(templateData.BodyHtmlTemplate, templateData.BodyModel)
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
