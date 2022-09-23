package handlers

import (
	"encoding/json"
	"net/http"

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
func RenderPdfFromHtmlHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &models.RenderData{}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		panic(err)
	}

	pdfService := pdf.NewPdfService(ctx)

	pdfData, err := pdfService.PdfFromHtml(data)

	if err != nil {
		panic(err)
	}

	if err := writePdf(ctx, w, pdfData); err != nil {
		panic(err)
	}
}
