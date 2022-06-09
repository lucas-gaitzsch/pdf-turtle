package templating

import "github.com/flosch/pongo2/v5"

const DjangoTemplateEngineKey = "django"

type DjangoTemplateEngine struct {
}

func (te *DjangoTemplateEngine) Execute(templateHtml *string, model interface{}) (*string, error) {
	empty := ""

	t, err := pongo2.FromString(*templateHtml)
	if err != nil {
		return &empty, err
	}

	html, err := t.Execute(pongo2.Context{"model": model})

	if err != nil {
		return &empty, err
	}

	return &html, nil
}

func (te *DjangoTemplateEngine) Test(templateHtml *string, model interface{}) error {
	_, err := te.Execute(templateHtml, model)

	return err
}
