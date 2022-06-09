package renderer

import (
	"bytes"
	"context"
	"io"
	"pdf-turtle/models"
	"pdf-turtle/templating"
	"pdf-turtle/utils"
	"pdf-turtle/utils/logging"
	"testing"
)

func TestRenderHtmlAsPdf(t *testing.T) {
	logging.InitTestLogger(t)
	defer logging.SetNullLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	html := "<b>test"
    renderer := NewAsyncHtmlRendererChromium(ctx, nil)
	defer renderer.Close()
	
	reader, err := utils.LogExecutionTimeWithResult("render pdf", nil, func() (io.Reader, error) {
		return renderer.RenderHtmlAsPdf(ctx, &models.RenderData{
			BodyHtml: &html,
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

    renderer := NewAsyncHtmlRendererChromium(ctx, nil)
	defer renderer.Close()

	reader, err := utils.LogExecutionTimeWithResult("render pdf", nil, func() (io.Reader, error) {
		return renderer.RenderHtmlAsPdf(ctx, &models.RenderData{
			BodyHtml: nil,
		})
	})

	if err != nil {
        t.Fatalf("RenderHtmlAsPdf with nil-pointer-body fails: %v", err)
    }

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	b := buf.Bytes()
	
	if len(b) == 0 || err != nil {
        t.Fatalf("RenderHtmlAsPdf with nil-pointer-body result empty; err: %v", err)
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
	
	htmlBody, err := utils.LogExecutionTimeWithResult("generate html from template", nil, func() (*string, error) {
		return templating.GetTemplateEngineByKey(templating.GoTemplateEngineKey).Execute(&template, data)
	})
	if err != nil {
        t.Fatalf("cant generate template %v", err)
    }

    renderer := NewAsyncHtmlRendererChromium(ctx, nil)
	defer renderer.Close()
	
	reader, err := utils.LogExecutionTimeWithResult("render pdf", nil, func() (io.Reader, error) {
		return renderer.RenderHtmlAsPdf(ctx, &models.RenderData{
			BodyHtml: htmlBody,
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

	return res;
}