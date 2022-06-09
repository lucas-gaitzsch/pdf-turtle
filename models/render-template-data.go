package models

import "encoding/json"

type RenderTemplateData struct {
	BodyHtmlTemplate   *string `json:"bodyHtmlTemplate"`
	HeaderHtmlTemplate string  `json:"headerHtmlTemplate,omitempty"`
	FooterHtmlTemplate string  `json:"footerHtmlTemplate,omitempty"`

	BodyModel   interface{} `json:"bodyModel,omitempty"`
	HeaderModel interface{} `json:"headerModel,omitempty"`
	FooterModel interface{} `json:"footerModel,omitempty"`

	TemplateEngine string `json:"templateEngine,omitempty" default:"golang" enums:"golang,handlebars,django"`

	RenderOptions RenderOptions `json:"options,omitempty"`
} // @name RenderTemplateData

//TODO:!! BodyHtmlTemplate -> example:"<b>Hello {{.name}}</b>"
//TODO:!! BodyModel -> swaggertype:"object,string" example:"name:Lucas"

func (data *RenderTemplateData) ParseJsonModelDataFromDoubleEncodedString() {
	data.BodyModel = parseJsonFieldFromDoubleEncodedString(data.BodyModel)
	data.HeaderModel = parseJsonFieldFromDoubleEncodedString(data.HeaderModel)
	data.FooterModel = parseJsonFieldFromDoubleEncodedString(data.FooterModel)
}

func parseJsonFieldFromDoubleEncodedString(model interface{}) interface{} {
	if str, ok := model.(string); ok {
		var parsed interface{}
		json.Unmarshal([]byte(str), &parsed)
		return parsed
	} else {
		return model
	}
}
