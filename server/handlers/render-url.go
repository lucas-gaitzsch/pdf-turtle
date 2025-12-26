package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lucas-gaitzsch/pdf-turtle/services/bundles"
	"github.com/lucas-gaitzsch/pdf-turtle/services/pdf"
)

type bytesOpener struct {
	data []byte
}

func (b *bytesOpener) Open() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(b.data)), nil
}

// RenderPdfFromUrlHandler
// GET endpoint that accepts URL as query parameter and renders PDF
// @Summary      Render PDF from URL
// @Description  Returns PDF file generated from HTML fetched from URL
// @Tags         Render URL
// @Accept       text/plain
// @Produce      application/pdf
// @Param        url  query  string  true  "URL to fetch HTML from"
// @Success      200         "PDF File"
// @Router       /api/pdf/from/url/render [get]
func RenderPdfFromUrlHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	urlParam := c.Query("url")
	if urlParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "url query parameter is required")
	}

	resp, err := http.Get(urlParam)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to fetch url: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bundle := &bundles.Bundle{}
	bundle.AddFile(bundles.BundleIndexFile, &bytesOpener{data: body})

	if err := bundle.TestIndexFile(); err != nil {
		return err
	}

	pdfService := pdf.NewPdfService(ctx)

	pdfData, err := pdfService.PdfFromBundle(bundle, "", "")
	if err != nil {
		return err
	}

	return writePdf(c, pdfData)
}
