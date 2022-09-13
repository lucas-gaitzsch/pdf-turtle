package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/loopback"
	"github.com/lucas-gaitzsch/pdf-turtle/models"

	"github.com/lucas-gaitzsch/pdf-turtle/services/bundleprovider"
	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
)

// RenderBundleHandler godoc
// @Summary      Render PDF from bundle and template provided in form-data (keys: bundle, model)
// @Description  Returns PDF file generated from HTML of body, header and footer
// @Tags         render html
// @Accept       multipart/form-data
// @Produce      application/pdf
// @Param        renderData  body  models.RenderData  true  "Render Data"
// @Success      200         "PDF File"
// @Router       /pdf/from/bundle/render [post]
func RenderBundleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conf := config.Get(ctx)

	r.ParseMultipartForm(int64(config.Get(ctx).MaxBodySizeInMb * 1024 * 1024))

	bundleFromForm, ok := r.MultipartForm.File["bundle"]

	if !ok || len(bundleFromForm) != 1 {
		//TODO:!! err
		return
	}

	reader, err := bundleFromForm[0].Open()

	bundle := bundleprovider.Bundle{}
	bundle.ReadFromZip(reader, bundleFromForm[0].Size)

	pdfService := pdf.NewPdfService(ctx)

	bundleProviderService := ctx.Value(config.ContextKeyBundleProviderService).(*bundleprovider.BundleProviderService)

	id, cleanup := bundleProviderService.Provide(bundle)
	defer cleanup()

	opt := bundle.GetOptions()
	opt.IsBundle = true
	opt.BasePath = fmt.Sprintf("http://127.0.0.1:%d%s/%s/", conf.LoopbackPort, loopback.BundlePath, id)

	var pdfData io.Reader
	var errRender error

	modelBody, hasTemplate := getValueFromForm(r.MultipartForm.Value, "model")
	if hasTemplate {
		templateEngine, _ := getValueFromForm(r.MultipartForm.Value, "templateEngine")

		templateData := &models.RenderTemplateData{
			HtmlTemplate:       bundle.GetBodyHtml(),
			HeaderHtmlTemplate: bundle.GetHeaderHtml(),
			FooterHtmlTemplate: bundle.GetFooterHtml(),
			TemplateEngine:     templateEngine,
			RenderOptions:      opt,
		}

		json.Unmarshal([]byte(modelBody), &templateData.Model)

		pdfData, errRender = pdfService.PdfFromHtmlTemplate(templateData)
	} else {
		data := &models.RenderData{
			Html:          bundle.GetBodyHtml(),
			HeaderHtml:    bundle.GetHeaderHtml(),
			FooterHtml:    bundle.GetFooterHtml(),
			RenderOptions: opt,
		}

		pdfData, errRender = pdfService.PdfFromHtml(data)
	}

	if errRender != nil {
		panic(err)
	}

	if err := writePdf(ctx, w, pdfData); err != nil {
		panic(err)
	}
}
