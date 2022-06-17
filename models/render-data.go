package models

import "github.com/lucas-gaitzsch/pdf-turtle/utils"

type RenderData struct {
	BodyHtml   *string `json:"bodyHtml" example:"<b>Hello World</b>"`
	HeaderHtml string  `json:"headerHtml,omitempty" example:"<h1>Heading</h1>"`
	FooterHtml string  `json:"footerHtml,omitempty" default:"<div class=\"default-footer\"><span class=\"pageNumber\"></span> of <span class=\"totalPages\"></span></div>"`

	RenderOptions RenderOptions `json:"options,omitempty"`
} // @name RenderData

func (d *RenderData) HasHeaderOrFooterHtml() bool {
	return d.HeaderHtml != "" || d.FooterHtml != ""
}

func (d *RenderData) SetDefaults() {
	utils.ReflectDefaultValues(d)
	
	d.RenderOptions.SetDefaults()
}
