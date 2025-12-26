package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/lucas-gaitzsch/pdf-turtle/config"
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
	conf := config.Get(ctx)

	urlParam := c.Query("url")
	if urlParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "url query parameter is required")
	}

	// Parse and validate URL
	parsedURL, err := url.Parse(urlParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid url format")
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fiber.NewError(fiber.StatusBadRequest, "only http and https schemes are allowed")
	}

	// Reconstruct URL to ensure it's properly formed
	safeURL := parsedURL.String()

	// Create HTTP client with optional proxy configuration
	client := &http.Client{}
	if proxyURL := conf.GetProxyUrl(); proxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Get(safeURL)
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
