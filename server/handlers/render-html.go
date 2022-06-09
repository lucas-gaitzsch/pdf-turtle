package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"pdf-turtle/models"
	"pdf-turtle/utils"
)

const TemplateEngineQueryKey = "template-engine"

// RenderPdfFromHtmlHandler godoc
// @Summary      Render PDF from HTML
// @Description  Returns PDF file generated from HTML of body, header and footer
// @Tags         render html
// @Accept       json
// @Produce      application/pdf
// @Param        renderData  body  models.RenderData  true  "Render Data"
// @Success      200         "PDF File"
// @Router       /pdf/from/html/render [post]
func RenderPdfFromHtmlHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pdfService := getPdfService(ctx)

	data := &models.RenderData{}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		panic(err)
	}

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
