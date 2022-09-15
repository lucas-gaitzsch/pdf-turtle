package templateengines

import (
	"encoding/json"
	"reflect"
	"testing"
)

const resultHtml = `
<html>
<body>
	<h1>Profile of Bruno</h1>
	<p>Working at Testcompany</p>
	<p>Locations:</p>
	<ul>
		<li>Chemnitz</li><li>Berlin</li><li>Amsterdam</li>
	</ul>
</body>
</html>
`

const jsonModel = `
{
	"name": "Bruno",
	"company": {
		"name": "Testcompany",
		"locations": ["Chemnitz", "Berlin", "Amsterdam"]
	}
}
`

func getModel() any {
	var model any
	json.Unmarshal([]byte(jsonModel), &model)
	return model
}

func TestGetTemplateEngineByKeyGo(t *testing.T) {
	templateStr := GoTemplateEngineKey

	engine, found := GetTemplateEngineByKey(HandlebarsTemplateEngineKey)

	if !found {
		t.Fatalf("cant find templateengine by key %s", templateStr)
	}

	if reflect.TypeOf(engine).Name() != reflect.TypeOf(&GoTemplateEngine{}).Name() {
		t.Fatalf("html not equal")
	}
}

func TestGetTemplateEngineByKeyHandlebars(t *testing.T) {
	templateStr := HandlebarsTemplateEngineKey

	engine, found := GetTemplateEngineByKey(templateStr)

	if !found {
		t.Fatalf("cant find templateengine by key %s", templateStr)
	}

	if reflect.TypeOf(engine).Name() != reflect.TypeOf(&HandlebarsTemplateEngine{}).Name() {
		t.Fatalf("html not equal")
	}
}

func TestGetTemplateEngineByKeyDjango(t *testing.T) {
	templateStr := DjangoTemplateEngineKey

	engine, found := GetTemplateEngineByKey(templateStr)

	if !found {
		t.Fatalf("cant find templateengine by key %s", templateStr)
	}

	if reflect.TypeOf(engine).Name() != reflect.TypeOf(&DjangoTemplateEngine{}).Name() {
		t.Fatalf("html not equal")
	}
}

func TestGetTemplateEngineByKeyEmpty(t *testing.T) {
	templateStr := ""

	engine, found := GetTemplateEngineByKey(templateStr)

	if found {
		t.Fatalf("empty key should not find any templateengine")
	}

	if reflect.TypeOf(engine).Name() != reflect.TypeOf(&DjangoTemplateEngine{}).Name() {
		t.Fatalf("html not equal")
	}
}

func TestGetTemplateEngineByKeyBullshit(t *testing.T) {
	templateStr := "assdfjkp"

	engine, found := GetTemplateEngineByKey(templateStr)

	if found {
		t.Fatalf("bullshit key should not find any templateengine")
	}

	if reflect.TypeOf(engine).Name() != reflect.TypeOf(&DjangoTemplateEngine{}).Name() {
		t.Fatalf("html not equal")
	}
}
