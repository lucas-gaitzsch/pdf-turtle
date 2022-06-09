package models

type RenderData struct {
	BodyHtml   *string `json:"bodyHtml" example:"<b>Hello World</b>"`
	HeaderHtml string  `json:"headerHtml,omitempty" example:"<h1>Heading</h1>"`
	FooterHtml string  `json:"footerHtml,omitempty" example:"<div class=\"default-footer\"><span class=\"pageNumber\"></span> of <span class=\"totalPages\"></span></div>"`

	RenderOptions RenderOptions `json:"options,omitempty"`
} // @name RenderData
