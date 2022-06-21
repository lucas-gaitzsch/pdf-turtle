package htmlparser

import (
	"strings"
	"testing"
)

func TestParseHeaderAndFooterAndCheckRemaining(t *testing.T) {
	p := New()

	doc := `
	<html>
		<PdfHeader>
			Header Content
		</PdfHeader>

		<PdfFooter>
			Footer Content
		</PdfFooter>

		<body>
			test
		</body>
	</html>`
	expectedDocRemaining := "<html><head></head><body>test</body></html>"

	p.Parse(&doc)

	header, footer := p.PopHeaderAndFooter()
	remaining, _ := p.GetHtml()

	if header != "Header Content" {
		t.Fatal("Header content was not parsed correctly")
	}

	if footer != "Footer Content" {
		t.Fatal("Footer content was not parsed correctly")
	}

	stripped := stripWhitespace(remaining)

	if stripped != expectedDocRemaining {
		t.Fatal("Header and Footer was not removed correctly")
	}
}

func TestParseOnlyHeader(t *testing.T) {
	p := New()

	doc := `
	<html>
		<PdfHeader>
			Header Content
		</PdfHeader>
	</html>`

	p.Parse(&doc)

	header, footer := p.PopHeaderAndFooter()

	if header != "Header Content" {
		t.Fatal("Header content was not parsed correctly")
	}

	if footer != "" {
		t.Fatal("Footer content should be empty")
	}
}

func TestAddStyleNoHead(t *testing.T) {
	p := New()

	doc := `
	<html>
		<body>test</body>
	</html>`

	shouldBe := `<html><head><style>body{color:red;}</style></head><body>test</body></html>`

	styles := "body{color:red;}"

	p.Parse(&doc)

	p.AddStyles(&styles)

	html, err := p.GetHtml()

	if err != nil {
		t.Fatalf("err should be nil: %v", err)
	}

	stripped := stripWhitespace(html)

	if stripped != shouldBe {
		t.Fatal("Style was not applied correctly")
	}
}

func TestAddStyleWithHead(t *testing.T) {
	p := New()

	doc := `
	<html>
		<head></head>
		<body>test</body>
	</html>`

	shouldBe := `<html><head><style>body{color:red;}</style></head><body>test</body></html>`

	styles := "body{color:red;}"

	p.Parse(&doc)

	p.AddStyles(&styles)

	html, err := p.GetHtml()

	if err != nil {
		t.Fatalf("err should be nil: %v", err)
	}

	stripped := stripWhitespace(html)

	if stripped != shouldBe {
		t.Fatal("Style was not applied correctly")
	}
}

func stripWhitespace(html *string) string {
	stripped := *html
	stripped = strings.ReplaceAll(stripped, " ", "")
	stripped = strings.ReplaceAll(stripped, "\n", "")
	stripped = strings.ReplaceAll(stripped, "\t", "")
	return stripped
}
