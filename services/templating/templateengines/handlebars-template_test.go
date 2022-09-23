package templateengines

import "testing"

const handlebarsTemplate = `
<html>
<body>
	<h1>Profile of {{name}}</h1>
	<p>Working at {{company.name}}</p>
	<p>Locations:</p>
	<ul>
		{{#each company.locations}}<li>{{this}}</li>{{/each}}
	</ul>
</body>
</html>
`

const handlebarsTemplateInvalid = `
<html>
<body>
	<h1>Profile of {{name}</h1>
	<p>Working at {{company.name}}</p>
	<p>Locations:</p>
	<ul>
		{{#each company.locations}}<li>{{this}}</li>{{/each}}
	</ul>
</body>
</html>
`

func TestHandlebarsTemplate(t *testing.T) {
	templateStr := handlebarsTemplate

	engine, _ := GetTemplateEngineByKey(HandlebarsTemplateEngineKey)

	htmlBody, err := engine.Execute(&templateStr, getModel())
	if err != nil {
		t.Fatalf("cant generate template %v", err)
	}

	if *htmlBody != resultHtml {
		t.Fatalf("html not equal")
	}
}

func TestHandlebarsTemplateTestValid(t *testing.T) {
	templateStr := handlebarsTemplate

	engine, _ := GetTemplateEngineByKey(HandlebarsTemplateEngineKey)

	err := engine.Test(&templateStr, getModel())
	if err != nil {
		t.Fatalf("cant generate template %v", err)
	}
}

func TestHandlebarsTemplateTestInvalid(t *testing.T) {
	templateStr := handlebarsTemplateInvalid

	engine, _ := GetTemplateEngineByKey(HandlebarsTemplateEngineKey)

	err := engine.Test(&templateStr, getModel())
	if err == nil {
		t.Fatalf("should fail")
	}
}
