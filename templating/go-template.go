package templating

import (
	"bytes"
	"html/template"
)

const GoTemplateEngineKey = "golang"

type GoTemplateEngine struct {
}

func (gte *GoTemplateEngine) Execute(templateHtml *string, model interface{}) (*string, error) {
	t, err := template.New("").Parse(*templateHtml)
	empty := ""
	if err != nil {
		return &empty, err
	}
	var buff bytes.Buffer

	if err := t.Execute(&buff, model); err != nil {
		return &empty, err
	}

	html := buff.String()

	return &html, nil
}

func (gte *GoTemplateEngine) Test(templateHtml *string, model interface{}) error {
	t, err := template.New("").Option("missingkey=error").Parse(*templateHtml)
	if err != nil {
		return err
	}

	var buff bytes.Buffer

	if err := t.Execute(&buff, model); err != nil {
		return err
	}

	return nil
}
