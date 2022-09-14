package templateengines

import (
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

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

func LogParsedTemplateEngine(keyToParse string, templateEngineInstance TemplateEngine,  found bool) {
	l := log.
		Debug().
		Str("givenTemplateEngine", keyToParse).
		Str("usedTemplateEngine", reflect.TypeOf(templateEngineInstance).String()).
		Bool("found", found)
	
	if keyToParse == "" {
		l.Msgf("template engine was not given -> fallback to %T", templateEngineInstance)			
	} else if found {
		l.Msg("given template engine was found")
	} else {
		l.Msg("given template engine was not found")
	}
}