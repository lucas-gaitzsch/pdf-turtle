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

	stripped := *remaining
	stripped = strings.ReplaceAll(stripped, " ", "")
	stripped = strings.ReplaceAll(stripped, "\n", "")
	stripped = strings.ReplaceAll(stripped, "\t", "")
	
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