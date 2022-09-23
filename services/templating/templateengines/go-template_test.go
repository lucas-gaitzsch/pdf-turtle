package templateengines

import (
	"testing"
)

const goTemplate = `
<html>
<body>
	<h1>Profile of {{.name}}</h1>
	<p>Working at {{.company.name}}</p>
	<p>Locations:</p>
	<ul>
		{{range $val := .company.locations}}<li>{{$val}}</li>{{end}}
	</ul>
</body>
</html>
`

const goTemplateUnknownProperty = `
<html>
<body>
	<h1>Profile of {{.lastname}}</h1>
	<p>Working at {{.company.name}}</p>
	<p>Locations:</p>
	<ul>
		{{range $val := .company.locations}}<li>{{$val}}</li>{{end}}
	</ul>
</body>
</html>
`

const goTemplateInvalid = `
<html>
<body>
	<h1>Profile of {{.name}</h1>
	<p>Working at {{.company.name}}</p>
	<p>Locations:</p>
	<ul>
		{{range $val := .company.locations}}<li>{{$val}}</li>{{end}}
	</ul>
</body>
</html>
`

func TestGoTemplate(t *testing.T) {
	templateStr := goTemplate

	engine, _ := GetTemplateEngineByKey(GoTemplateEngineKey)

	htmlBody, err := engine.Execute(&templateStr, getModel())
	if err != nil {
		t.Fatalf("cant generate template %v", err)
	}

	if *htmlBody != resultHtml {
		t.Fatalf("html not equal")
	}
}

func TestGoTemplateTestValid(t *testing.T) {
	templateStr := goTemplate

	engine, _ := GetTemplateEngineByKey(GoTemplateEngineKey)

	err := engine.Test(&templateStr, getModel())
	if err != nil {
		t.Fatalf("cant generate template %v", err)
	}
}

func TestGoTemplateTestInvalid(t *testing.T) {
	templateStr := goTemplateInvalid

	engine, _ := GetTemplateEngineByKey(GoTemplateEngineKey)

	err := engine.Test(&templateStr, getModel())
	if err == nil {
		t.Fatalf("should fail")
	}
}

func TestGoTemplateTestUnknownProperty(t *testing.T) {
	templateStr := goTemplateUnknownProperty

	engine, _ := GetTemplateEngineByKey(GoTemplateEngineKey)

	err := engine.Test(&templateStr, getModel())
	if err == nil {
		t.Fatalf("should fail")
	}
}
