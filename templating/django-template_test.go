package templating

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
	
	htmlBody, err := GetTemplateEngineByKey(DjangoTemplateEngineKey).Execute(&templateStr, getModel())
	if err != nil {
        t.Fatalf("cant generate template %v", err)
    }

	if *htmlBody != resultHtml {
        t.Fatalf("html not equal")
    }
}

func TestDjangoTemplateTestValid(t *testing.T) {
	templateStr := djangoTemplate
	
	err := GetTemplateEngineByKey(DjangoTemplateEngineKey).Test(&templateStr, getModel())
	if err != nil {
        t.Fatalf("cant generate template %v", err)
    }
}

func TestDjangoTemplateTestInvalid(t *testing.T) {
	templateStr := djangoTemplateInvalid
	
	err := GetTemplateEngineByKey(DjangoTemplateEngineKey).Test(&templateStr, getModel())
	if err == nil {
        t.Fatalf("should fail")
    }
}
