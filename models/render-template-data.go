package models

import "encoding/json"

type RenderTemplateData struct {
	HtmlTemplate       *string `json:"htmlTemplate"`
	HeaderHtmlTemplate string  `json:"headerHtmlTemplate,omitempty"` // Optional template for header. If empty, the header template will be parsed from main template (<PdfHeader></PdfHeader>).
	FooterHtmlTemplate string  `json:"footerHtmlTemplate,omitempty"` // Optional template for footer. If empty, the footer template will be parsed from main template (<PdfFooter></PdfFooter>).

	Model       any `json:"model,omitempty"`
	HeaderModel any `json:"headerModel,omitempty"` // Optional model for header. If empty or null model was used.
	FooterModel any `json:"footerModel,omitempty"` // Optional model for footer. If empty or null model was used.

	TemplateEngine string `json:"templateEngine,omitempty" default:"golang" enums:"golang,handlebars,django"`

	RenderOptions RenderOptions `json:"options,omitempty"`
} // @name RenderTemplateData

//TODO:!! HtmlTemplate -> example:"<b>Hello {{.name}}</b>"
//TODO:!! BodyModel -> swaggertype:"object,string" example:"name:Lucas"

func (d *RenderTemplateData) HasHeaderOrFooterHtml() bool {
	return d.HeaderHtmlTemplate != "" || d.FooterHtmlTemplate != ""
}

func (d *RenderTemplateData) GetBodyHtml() *string {
	return d.HtmlTemplate
}

func (d *RenderTemplateData) SetBodyHtml(html *string) {
	d.HtmlTemplate = html
}

func (d *RenderTemplateData) GetHeaderHtml() string {
	return d.HeaderHtmlTemplate
}

func (d *RenderTemplateData) SetHeaderHtml(html string) {
	d.HeaderHtmlTemplate = html
}

func (d *RenderTemplateData) GetFooterHtml() string {
	return d.FooterHtmlTemplate
}

func (d *RenderTemplateData) SetFooterHtml(html string) {
	d.FooterHtmlTemplate = html
}

func (d *RenderTemplateData) HasBuiltinStylesExcluded() bool {
	return d.RenderOptions.ExcludeBuiltinStyles
}

func (d *RenderTemplateData) GetBodyModel() any {
	return d.Model
}

func (d *RenderTemplateData) HasHeaderOrFooterModel() bool {
	return d.HeaderModel != nil || d.FooterModel != nil
}

func (d *RenderTemplateData) GetHeaderModel() any {
	if !d.HasHeaderOrFooterModel() {
		return d.Model
	}
	return d.HeaderModel
}

func (d *RenderTemplateData) GetFooterModel() any {
	if !d.HasHeaderOrFooterModel() {
		return d.Model
	}
	return d.FooterModel
}

func (d *RenderTemplateData) ParseJsonModelDataFromDoubleEncodedString() {
	d.Model = parseJsonFieldFromDoubleEncodedString(d.Model)
	d.HeaderModel = parseJsonFieldFromDoubleEncodedString(d.HeaderModel)
	d.FooterModel = parseJsonFieldFromDoubleEncodedString(d.FooterModel)
}
func parseJsonFieldFromDoubleEncodedString(model any) any {
	if str, ok := model.(string); ok {
		var parsed any
		json.Unmarshal([]byte(str), &parsed)
		return parsed
	} else {
		return model
	}
}
