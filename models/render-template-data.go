package models

import "encoding/json"

type RenderTemplateData struct {
	HtmlTemplate       *string `json:"htmlTemplate"`
	// Optional template for header. If empty, the header template will be parsed from main template (<PdfHeader></PdfHeader>).
	HeaderHtmlTemplate string  `json:"headerHtmlTemplate,omitempty"` 
	// Optional template for footer. If empty, the footer template will be parsed from main template (<PdfFooter></PdfFooter>).
	FooterHtmlTemplate string  `json:"footerHtmlTemplate,omitempty"` 

	// Model with your data matching to the templates
	Model any `json:"model,omitempty" swaggertype:"object"`

	TemplateEngine string `json:"templateEngine,omitempty" default:"golang" enums:"golang,handlebars,django"`

	RenderOptions RenderOptions `json:"options,omitempty"`
} // @name RenderTemplateData

func (d *RenderTemplateData) HasHeaderOrFooterHtml() bool {
	return d.HeaderHtmlTemplate != "" || d.FooterHtmlTemplate != ""
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

func (d *RenderTemplateData) ParseJsonModelDataFromDoubleEncodedString() {
	d.Model = parseJsonFieldFromDoubleEncodedString(d.Model)
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
