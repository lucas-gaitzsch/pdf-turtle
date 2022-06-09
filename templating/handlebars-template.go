package templating

import "github.com/aymerick/raymond"

const HandlebarsTemplateEngineKey = "handlebars"

type HandlebarsTemplateEngine struct {
}

func (te *HandlebarsTemplateEngine) Execute(templateHtml *string, model interface{}) (*string, error) {
	empty := ""

	html, err := raymond.Render(*templateHtml, model)

	if err != nil {
		return &empty, err
	}

	return &html, nil
}

func (te *HandlebarsTemplateEngine) Test(templateHtml *string, model interface{}) error {
	_, err := te.Execute(templateHtml, model)
	
	return err
}
