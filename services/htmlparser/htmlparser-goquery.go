package htmlparser

import (
	"strings"

	"github.com/lucas-gaitzsch/pdf-turtle/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

func (p *HtmlParserGoQuery) AddStyles(css *string) {
	headTag := p.doc.Find("head").First()
	if headTag == nil {
		log.Panic().Msg("html should have a head")
		return
	}

	styleNode := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Style,
		Data:     "style",
	}

	styleNode.AppendChild(&html.Node{
		Type: html.RawNode,
		Data: *css,
	})

	headTag.Nodes[0].AppendChild(styleNode)
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
		header = utils.TrimStrWhitespace(html)
		headerNode.Remove()
	}

	footerNode := p.doc.Find(FooterNodeTag).First()
	if footerNode != nil {
		html, _ := footerNode.Html()
		footer = utils.TrimStrWhitespace(html)
		footerNode.Remove()
	}

	return
}

func (p *HtmlParserGoQuery) GetHtml() (*string, error) {
	html, err := p.doc.Html()
	return &html, err
}
