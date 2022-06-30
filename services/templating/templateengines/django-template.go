package templateengines

import "github.com/flosch/pongo2/v5"

const DjangoTemplateEngineKey = "django"

type DjangoTemplateEngine struct {
}

func (te *DjangoTemplateEngine) Execute(templateHtml *string, model any) (*string, error) {
	empty := ""

	t, err := pongo2.FromString(*templateHtml)
	if err != nil {
		return &empty, err
	}

	html, err := t.Execute(pongo2.Context{
		"model": model,
		"func":  templateFunctions,
	})

	if err != nil {
		return &empty, err
	}

	return &html, nil
}

func (te *DjangoTemplateEngine) Test(templateHtml *string, model any) error {
	_, err := te.Execute(templateHtml, model)

	return err
}
