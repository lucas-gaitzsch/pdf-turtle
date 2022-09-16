package templateengines

import (
	"errors"

	"github.com/aymerick/raymond"
)

const HandlebarsTemplateEngineKey = "handlebars"

type HandlebarsTemplateEngine struct {
}

func (te *HandlebarsTemplateEngine) Execute(templateHtml *string, model any) (*string, error) {
	empty := ""

	if templateHtml == nil {
		return &empty, errors.New("templateHtml is nil")
	}

	t, err := raymond.Parse(*templateHtml)
	if err != nil {
		return &empty, err
	}

	t.RegisterHelpers(templateFunctions)

	html, err := t.Exec(model)
	if err != nil {
		return &empty, err
	}

	return &html, nil
}

func (te *HandlebarsTemplateEngine) Test(templateHtml *string, model any) error {
	_, err := te.Execute(templateHtml, model)

	return err
}
