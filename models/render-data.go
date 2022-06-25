package models

import "github.com/lucas-gaitzsch/pdf-turtle/utils"

type RenderData struct {
	Html       *string `json:"html" example:"<b>Hello World</b>"`
	HeaderHtml string  `json:"headerHtml,omitempty" example:"<h1>Heading</h1>"`                                                                                             // Optional html for header. If empty, the header html will be parsed from main html (<PdfHeader></PdfHeader>).
	FooterHtml string  `json:"footerHtml,omitempty" default:"<div class=\"default-footer\"><div><span class=\"pageNumber\"></span> of <span class=\"totalPages\"></span></div></div>"` // Optional html for footer. If empty, the footer html will be parsed from main html (<PdfFooter></PdfFooter>).

	RenderOptions RenderOptions `json:"options,omitempty"`
} // @name RenderData

func (d *RenderData) HasHeaderOrFooterHtml() bool {
	return d.HeaderHtml != "" || d.FooterHtml != ""
}

func (d *RenderData) GetBodyHtml() *string {
	return d.Html
}

func (d *RenderData) SetBodyHtml(html *string) {
	d.Html = html
}

func (d *RenderData) GetHeaderHtml() string {
	return d.HeaderHtml
}

func (d *RenderData) SetHeaderHtml(html string) {
	d.HeaderHtml = html
}

func (d *RenderData) GetFooterHtml() string {
	return d.FooterHtml
}

func (d *RenderData) SetFooterHtml(html string) {
	d.FooterHtml = html
}

func (d *RenderData) HasBuiltinStylesExcluded() bool {
	return d.RenderOptions.ExcludeBuiltinStyles
}

func (d *RenderData) SetDefaults() {
	utils.ReflectDefaultValues(d)

	d.RenderOptions.SetDefaults()
}
