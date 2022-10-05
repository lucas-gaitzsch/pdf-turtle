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
	engine:= getTemplateEngineByKey(t, GoTemplateEngineKey)

	if reflect.TypeOf(engine).Elem().Name() != reflect.TypeOf(GoTemplateEngine{}).Name() {
		fatalWrongTemplateEngine(t)
	}
}

func TestGetTemplateEngineByKeyHandlebars(t *testing.T) {
	engine:= getTemplateEngineByKey(t, HandlebarsTemplateEngineKey)

	if reflect.TypeOf(engine).Elem().Name() != reflect.TypeOf(HandlebarsTemplateEngine{}).Name() {
		fatalWrongTemplateEngine(t)
	}
}

func TestGetTemplateEngineByKeyDjango(t *testing.T) {
	engine:= getTemplateEngineByKey(t, DjangoTemplateEngineKey)

	if reflect.TypeOf(engine).Elem().Name() != reflect.TypeOf(DjangoTemplateEngine{}).Name() {
		fatalWrongTemplateEngine(t)
	}
}

func TestGetTemplateEngineByKeyEmpty(t *testing.T) {
	engine, found := GetTemplateEngineByKey("")

	if found {
		t.Fatal("found templateengine by empty key")
	}


	if reflect.TypeOf(engine).Elem().Name() != reflect.TypeOf(GoTemplateEngine{}).Name() {
		fatalWrongTemplateEngine(t)
	}
}

func TestGetTemplateEngineByKeyBullshit(t *testing.T) {
	engine, found := GetTemplateEngineByKey("bullshit")

	if found {
		t.Fatalf("bullshit key should not find any templateengine")
	}

	if reflect.TypeOf(engine).Elem().Name() != reflect.TypeOf(GoTemplateEngine{}).Name() {
		fatalWrongTemplateEngine(t)
	}
}

func getTemplateEngineByKey(t *testing.T, key string) TemplateEngine{
	engine, found := GetTemplateEngineByKey(key)

	if !found {
		t.Fatalf("cant find templateengine by key '%s'", key)
	}

	return engine
}

func fatalWrongTemplateEngine(t *testing.T) {
	t.Fatal("wrong template engine was loaded")
}