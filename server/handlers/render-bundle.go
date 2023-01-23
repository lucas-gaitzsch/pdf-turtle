package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/lucas-gaitzsch/pdf-turtle/config"

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
// @Description  Returns PDF file generated from bundle (Zip-File) of HTML or HTML template of body, header, footer and assets. The index.html file in the Zip-Bundle is required
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

	r.ParseMultipartForm(int64(config.Get(ctx).MaxBodySizeInMb) * 1024 * 1024)

	bundleFromForm, ok := r.MultipartForm.File[formDataKeyBundle]

	if !ok || len(bundleFromForm) == 0 {
		panic(errors.New("no bundle data with key 'bundle' was attached in form data"))
	}

	bundle := &bundles.Bundle{}

	for _, fb := range bundleFromForm {

		if strings.HasPrefix(fb.Filename, "bundle") || strings.HasPrefix(fb.Filename, "blob") || strings.HasSuffix(fb.Filename, ".zip") {
			reader, err := fb.Open()
			if err != nil {
				panic(err)
			}
			defer reader.Close()

			err = bundle.ReadFromZip(reader, fb.Size)

			if err != nil {
				panic(err)
			}
		} else {
			fp := &bundles.OpenerFileProxy{
				MultipartFileOpener: fb,
			}
			bundle.AddFile(fb.Filename, fp)
		}
	}

	err := bundle.TestIndexFile()
	if err != nil {
		panic(err)
	}

	pdfService := pdf.NewPdfService(ctx)

	jsonModel, _ := getValueFromForm(r.MultipartForm.Value, formDataKeyModel)
	templateEngine, _ := getValueFromForm(r.MultipartForm.Value, formDataKeyTemplateEngine)

	pdfData, errRender := pdfService.PdfFromBundle(bundle, jsonModel, templateEngine)

	if errRender != nil {
		panic(errRender)
	}

	if err := writePdf(ctx, w, pdfData); err != nil {
		panic(err)
	}
}
