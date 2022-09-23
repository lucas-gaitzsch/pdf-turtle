package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/loopback"
	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/rs/zerolog/log"

	"github.com/lucas-gaitzsch/pdf-turtle/services/bundles"
	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
)

const (
	formDataKeyBundle         = "bundle"
	formDataKeyModel          = "model"
	formDataKeyTemplateEngine = "templateEngine"
)

// RenderBundleHandler godoc
// @Summary      Render PDF from bundle including HTML(-Template) with model and assets provided in form-data (keys: bundle, model)
// @Description  Returns PDF file generated from bundle (Zip-File) of HTML or HTML template of body, header, footer and assets. The index.html file in the Zip-Bundle is required.
// @Tags         Render HTML-Bundle
// @Accept       multipart/form-data
// @Produce      application/pdf
// @Param        bundle          formData  file    true   "Bundle Zip-File"
// @Param        model           formData  string  false  "JSON-Model for template (only required for template)"
// @Param        templateEngine  formData  string  false  "Template engine to use for template (only required for template)"
// @Success      200             "PDF File"
// @Router       /api/pdf/from/html-bundle/render [post]
func RenderBundleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conf := config.Get(ctx)

	r.ParseMultipartForm(int64(config.Get(ctx).MaxBodySizeInMb) * 1024 * 1024)

	bundleFromForm, ok := r.MultipartForm.File[formDataKeyBundle]

	if !ok || len(bundleFromForm) == 0 {
		panic(errors.New("no zip bundle with key 'bundle' was attached in form data"))
	}

	bundle := bundles.Bundle{}

	for _, fb := range bundleFromForm {
		reader, err := fb.Open()
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		err = bundle.ReadFromZip(reader, fb.Size)

		if err != nil {
			panic(err)
		}
	}

	pdfService := pdf.NewPdfService(ctx)

	bundleProviderService := ctx.Value(config.ContextKeyBundleProviderService).(*bundles.BundleProviderService)

	id, cleanup := bundleProviderService.Provide(bundle)
	defer cleanup()

	opt := bundle.GetOptions()
	opt.BasePath = fmt.Sprintf("http://127.0.0.1:%d%s/%s/", conf.LoopbackPort, loopback.BundlePath, id)

	var pdfData io.Reader
	var errRender error

	modelBody, hasModel := getValueFromForm(r.MultipartForm.Value, formDataKeyModel)
	hasModelLoggingPreparation := log.Debug().Bool("hasModel", hasModel)
	if hasModel {
		hasModelLoggingPreparation.Msg("got model in form data -> render with template engine")

		templateEngine, hasTemplateEngine := getValueFromForm(r.MultipartForm.Value, formDataKeyTemplateEngine)
		if hasTemplateEngine {
			log.Debug().
				Str("templateEngine", templateEngine).
				Msg("got templateEngine in form data")
		}

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
		hasModelLoggingPreparation.Msg("no model given with key 'model' in form data -> render plain html")

		data := &models.RenderData{
			Html:          bundle.GetBodyHtml(),
			HeaderHtml:    bundle.GetHeaderHtml(),
			FooterHtml:    bundle.GetFooterHtml(),
			RenderOptions: opt,
		}

		pdfData, errRender = pdfService.PdfFromHtml(data)
	}

	if errRender != nil {
		panic(errRender)
	}

	if err := writePdf(ctx, w, pdfData); err != nil {
		panic(err)
	}
}
