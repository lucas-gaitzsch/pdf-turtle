package templateengines

import "strings"

type TemplateEngine interface {
	Execute(templateHtml *string, model any) (*string, error)
	Test(templateHtml *string, model any) error
}

func GetTemplateEngineByKey(key string) (TemplateEngine, bool) {
	var templateEngine TemplateEngine

	switch strings.ToLower(key) {
	case strings.ToLower(HandlebarsTemplateEngineKey):
		templateEngine = &HandlebarsTemplateEngine{}
	case strings.ToLower(DjangoTemplateEngineKey):
		templateEngine = &DjangoTemplateEngine{}
	case strings.ToLower(GoTemplateEngineKey):
		templateEngine = &GoTemplateEngine{}
	default:
		return &GoTemplateEngine{}, false
	}

	return templateEngine, true
}
