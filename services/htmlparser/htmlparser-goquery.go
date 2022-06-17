package htmlparser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

type HtmlParserGoQuery struct {
	doc *goquery.Document
}

func (p *HtmlParserGoQuery) Parse(document *string) error {
	r := strings.NewReader(*document)

	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		log.Error().Err(err).Msg("cant parse dom from body html")
		return err
	}

	p.doc = doc

	return nil
}

func (p *HtmlParserGoQuery) PopHeaderAndFooter() (header string, footer string) {
	if p.doc == nil {
		log.Panic().Msg("parsedDoc==nil -> please call .Parse(doc) first")
	}

	header = ""
	footer = ""
	
	headerNode := p.doc.Find(HeaderNodeTag).First()
	if headerNode != nil {
		html, _ := headerNode.Html()
		header = trim(html)
		headerNode.Remove()
	}

	footerNode := p.doc.Find(FooterNodeTag).First()
	if footerNode != nil {
		html, _ := footerNode.Html()
		footer = trim(html)
		footerNode.Remove()
	}

	return
}


func (p *HtmlParserGoQuery) GetHtml() (*string, error) {
	html, err := p.doc.Html()
	return &html, err
}

func trim(str string) string {
	return strings.Trim(str, trimCutset)
}