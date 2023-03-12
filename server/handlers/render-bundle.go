package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

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
func RenderBundleHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	bundlesFromForm, ok := form.File[formDataKeyBundle]
	if !ok || len(bundlesFromForm) == 0 {
		return errors.New("no bundle data with key 'bundle' was attached in form data")
	}

	bundle := &bundles.Bundle{}

	for _, fb := range bundlesFromForm {

		if strings.HasPrefix(fb.Filename, "bundle") || fb.Header.Get("Content-Type") == "application/zip" || strings.HasSuffix(fb.Filename, ".zip") {
			reader, err := fb.Open()
			if err != nil {
				return err
			}
			defer reader.Close()

			err = bundle.ReadFromZip(reader, fb.Size)

			if err != nil {
				return err
			}
		} else {
			fp := &bundles.OpenerFileProxy{
				MultipartFileOpener: fb,
			}
			bundle.AddFile(fb.Filename, fp)
		}
	}

	err = bundle.TestIndexFile()
	if err != nil {
		return err
	}

	pdfService := pdf.NewPdfService(ctx)

	jsonModel, _ := getValueFromForm(form.Value, formDataKeyModel)
	templateEngine, _ := getValueFromForm(form.Value, formDataKeyTemplateEngine)

	pdfData, errRender := pdfService.PdfFromBundle(bundle, jsonModel, templateEngine)

	if errRender != nil {
		return errRender
	}

	return writePdf(c, pdfData)
}
