package utils

import (
	"reflect"
	"testing"
)

const fatalMsgDefaultNotAsExpected = "struct defaults was not set as expected"

type testStruct struct {
	TestProp1 string  `default:"testteststr"`
	TestProp2 int     `default:"5"`
	TestProp3 bool    `default:"true"`
	TestProp4 *bool   `default:"true"`
	TestProp5 *string `default:"testteststr2"`
	TestProp6 string
	testProp7 bool `default:"true"`
}

func TestReflectDefaultValuesEmptyStruct(t *testing.T) {
	testBool := true
	testStr := "testteststr2"
	shouldBe := testStruct{
		TestProp1: "testteststr",
		TestProp2: 5,
		TestProp3: true,
		TestProp4: &testBool,
		TestProp5: &testStr,
	}

	s := &testStruct{}
	ReflectDefaultValues(s)

	if !reflect.DeepEqual(shouldBe, *s) {
		t.Fatal(fatalMsgDefaultNotAsExpected)
	}
}

func TestReflectDefaultValuesPartiallyPrefilledStruct(t *testing.T) {
	testBool := true
	testStr := "peter"
	shouldBe := testStruct{
		TestProp1: "peter",
		TestProp2: 5,
		TestProp3: true,
		TestProp4: &testBool,
		TestProp5: &testStr,
	}

	s := &testStruct{
		TestProp1: "peter",
		TestProp5: &testStr,
	}
	ReflectDefaultValues(s)

	if !reflect.DeepEqual(shouldBe, *s) {
		t.Fatal(fatalMsgDefaultNotAsExpected)
	}
}

func TestReflectDefaultValuesNoDefaultAnnotation(t *testing.T) {

	type testStructNoDefaults struct {
		TestProp1 string
		TestProp2 int
		TestProp3 bool
		TestProp4 *bool
		TestProp5 *string
		TestProp6 string
		testProp7 bool
	}

	shouldBe := testStructNoDefaults{
		TestProp1: "peter",
	}

	s := &testStructNoDefaults{
		TestProp1: "peter",
	}
	ReflectDefaultValues(s)

	if !reflect.DeepEqual(shouldBe, *s) {
		t.Fatal(fatalMsgDefaultNotAsExpected)
	}
}
