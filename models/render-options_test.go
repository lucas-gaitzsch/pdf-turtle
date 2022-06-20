package models

import (
	"reflect"
	"testing"
)

const fatalMsgDefaultNotAsExpected = "struct defaults was not set as expected"

func TestSetDefaultsEmptyStruct(t *testing.T) {
	shouldBe := &RenderOptions{
		PageFormat: "A4",
		PageSize: PageSize{
			Width:  210,
			Height: 297,
		},
		Landscape:            false,
		ExcludeBuiltinStyles: false,
		Margins: &RenderOptionsMargins{
			Top:    25,
			Right:  25,
			Bottom: 20,
			Left:   25,
		},
	}

	testObj := &RenderOptions{}
	testObj.SetDefaults()

	if !reflect.DeepEqual(shouldBe, testObj) {
		t.Fatal(fatalMsgDefaultNotAsExpected)
	}
}

func TestSetDefaultsPrefilledStruct(t *testing.T) {
	shouldBe := &RenderOptions{
		PageFormat: "A4",
		PageSize: PageSize{
			Width:  210,
			Height: 297,
		},
		Landscape:            true,
		ExcludeBuiltinStyles: false,
		Margins: &RenderOptionsMargins{
			Top:    25,
			Right:  25,
			Bottom: 20,
			Left:   25,
		},
	}

	testObj := &RenderOptions{
		Landscape: true,
	}
	testObj.SetDefaults()

	if !reflect.DeepEqual(shouldBe, testObj) {
		t.Fatal(fatalMsgDefaultNotAsExpected)
	}
}

func TestSetDefaultsWithFormat(t *testing.T) {
	shouldBe := &RenderOptions{
		PageFormat: "A3",
		PageSize: PageSize{
			Width:  297,
			Height: 420,
		},
		Landscape:            false,
		ExcludeBuiltinStyles: false,
		Margins: &RenderOptionsMargins{
			Top:    25,
			Right:  25,
			Bottom: 20,
			Left:   25,
		},
	}

	testObj := &RenderOptions{
		PageFormat: "A3",
	}
	testObj.SetDefaults()

	if !reflect.DeepEqual(shouldBe, testObj) {
		t.Fatal(fatalMsgDefaultNotAsExpected)
	}
}

func TestSetDefaultsWithPageSize(t *testing.T) {
	shouldBe := &RenderOptions{
		PageFormat: "A4",
		PageSize: PageSize{
			Width:  200,
			Height: 300,
		},
		Landscape:            false,
		ExcludeBuiltinStyles: false,
		Margins: &RenderOptionsMargins{
			Top:    25,
			Right:  25,
			Bottom: 20,
			Left:   25,
		},
	}

	testObj := &RenderOptions{
		PageSize: PageSize{
			Width:  200,
			Height: 300,
		},
	}
	testObj.SetDefaults()

	if !reflect.DeepEqual(shouldBe, testObj) {
		t.Fatal(fatalMsgDefaultNotAsExpected)
	}
}
