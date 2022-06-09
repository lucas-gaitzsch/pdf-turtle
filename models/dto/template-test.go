package dto

type TemplateTestResult struct {
	IsValid             bool    `json:"isValid"`
	BodyTemplateError   *string `json:"bodyTemplateError"`
	HeaderTemplateError *string `json:"headerTemplateError"`
	FooterTemplateError *string `json:"footerTemplateError"`
} // @name TemplateTestResult
