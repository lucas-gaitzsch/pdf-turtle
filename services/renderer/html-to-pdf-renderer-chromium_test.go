package renderer

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services/templating/templateengines"
	"github.com/lucas-gaitzsch/pdf-turtle/utils/logging"
)

func TestRenderHtmlAsPdf(t *testing.T) {
	logging.InitTestLogger(t)
	defer logging.SetNullLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	html := "<b>test"
	renderer := NewAsyncHtmlRendererChromium(ctx)
	defer renderer.Close()

	reader, err := logging.LogExecutionTimeWithResults("render pdf", nil, func() (io.Reader, error) {
		return renderer.RenderHtmlAsPdf(ctx, &models.RenderData{
			Html:       &html,
			HeaderHtml: html,
			FooterHtml: html,
		})
	})
	if err != nil {
		t.Fatalf("RenderHtmlAsPdf fails: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	b := buf.Bytes()

	if len(b) == 0 || err != nil {
		t.Fatalf("RenderHtmlAsPdf result empty; err: %v", err)
	}
}

func TestRenderHtmlAsPdfWithNilPointerBody(t *testing.T) {
	logging.InitTestLogger(t)
	defer logging.SetNullLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	renderer := NewAsyncHtmlRendererChromium(ctx)
	defer renderer.Close()

	reader, err := logging.LogExecutionTimeWithResults("render pdf (nil body)", nil, func() (io.Reader, error) {
		return renderer.RenderHtmlAsPdf(ctx, &models.RenderData{
			Html: nil,
		})
	})

	if err == nil {
		t.Fatalf("RenderHtmlAsPdf with nil-pointer-body should fail: %v", err)
	}

	if reader != nil {
		t.Fatal("Reader should be nil")
	}
}

func TestRenderHugeHtmlAsPdf(t *testing.T) {
	logging.InitTestLogger(t)
	defer logging.SetNullLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data := generateRange(10_000)
	template := `
	<table>
		{{range $val := .}}
		<tr><td>{{$val}}</td></tr>
		{{end}}
	</table>`

	htmlBody, err := logging.LogExecutionTimeWithResults("generate html from template", nil, func() (*string, error) {
		engine, _ := templateengines.GetTemplateEngineByKey(templateengines.GoTemplateEngineKey)
		return engine.Execute(&template, data)
	})
	if err != nil {
		t.Fatalf("cant generate template %v", err)
	}

	renderer := NewAsyncHtmlRendererChromium(ctx)
	defer renderer.Close()

	reader, err := logging.LogExecutionTimeWithResults("render pdf (huge)", nil, func() (io.Reader, error) {
		return renderer.RenderHtmlAsPdf(ctx, &models.RenderData{
			Html:       htmlBody,
			HeaderHtml: "<h1 id=\"header-template\" style=\"font-size:3mm !important;\">Heading</h1>",
			RenderOptions: models.RenderOptions{
				Margins: &models.RenderOptionsMargins{
					Top: 10,
				},
			},
		})
	})
	if err != nil {
		t.Fatalf("RenderHtmlAsPdf fails: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	b := buf.Bytes()

	if len(b) == 0 || err != nil {
		t.Fatalf("RenderHtmlAsPdf result empty; err: %v", err)
	}
}

func generateRange(until int) []int {
	res := make([]int, until)

	for i := 0; i < until; i++ {
		res[i] = i
	}

	return res
}
