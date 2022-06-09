package templating

import "strings"

type TemplateEngine interface {
	Execute(templateHtml *string, model interface{}) (*string, error)
	Test(templateHtml *string, model interface{}) error
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
