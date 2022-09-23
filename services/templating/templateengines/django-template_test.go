package templateengines

import (
	"testing"
)

const djangoTemplate = `
<html>
<body>
	<h1>Profile of {{model.name}}</h1>
	<p>Working at {{model.company.name}}</p>
	<p>Locations:</p>
	<ul>
		{% for loc in model.company.locations %}<li>{{ loc }}</li>{% endfor %}
	</ul>
</body>
</html>
`

const djangoTemplateInvalid = `
<html>
<body>
	<h1>Profile of {{model.name}</h1>
	<p>Working at {{model.company.name}}</p>
	<p>Locations:</p>
	<ul>
		{% for loc in model.company.locations %}<li>{{ loc }}</li>{% endfor %}
	</ul>
</body>
</html>
`

func TestDjangoTemplate(t *testing.T) {
	templateStr := djangoTemplate

	engine, _ := GetTemplateEngineByKey(DjangoTemplateEngineKey)

	htmlBody, err := engine.Execute(&templateStr, getModel())
	if err != nil {
		t.Fatalf("cant generate template %v", err)
	}

	if *htmlBody != resultHtml {
		t.Fatalf("html not equal")
	}
}

func TestDjangoTemplateTestValid(t *testing.T) {
	templateStr := djangoTemplate

	engine, _ := GetTemplateEngineByKey(DjangoTemplateEngineKey)

	err := engine.Test(&templateStr, getModel())
	if err != nil {
		t.Fatalf("cant generate template %v", err)
	}
}

func TestDjangoTemplateTestInvalid(t *testing.T) {
	templateStr := djangoTemplateInvalid

	engine, _ := GetTemplateEngineByKey(DjangoTemplateEngineKey)

	err := engine.Test(&templateStr, getModel())
	if err == nil {
		t.Fatalf("should fail")
	}
}
