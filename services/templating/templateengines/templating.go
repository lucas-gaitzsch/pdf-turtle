package templateengines

import "strings"

type TemplateEngine interface {
	Execute(templateHtml *string, model any) (*string, error)
	Test(templateHtml *string, model any) error
}

func GetTemplateEngineByKey(key string) TemplateEngine {
	switch strings.ToLower(key) {
	case strings.ToLower(HandlebarsTemplateEngineKey):
		return &HandlebarsTemplateEngine{}
	case strings.ToLower(DjangoTemplateEngineKey):
		return &DjangoTemplateEngine{}
	case strings.ToLower(GoTemplateEngineKey):
		return &GoTemplateEngine{}
	default:
		return &GoTemplateEngine{}
	}
}
