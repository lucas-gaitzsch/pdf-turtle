package utils

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestAppendStyleToHtml(t *testing.T) {
	inputHtml := "<body>Hello World</body>"
	css := ".test { color: red; }"
	outputPtr := AppendStyleToHtml(&inputHtml, &css)

	if outputPtr == nil {
		t.Fatal("Result should not be nil")
	}

	output := *outputPtr

	if !strings.HasSuffix(output, "<style>" + css + "</style>") {
		t.Fatal("Result should have appended css")
	}
	if !strings.HasPrefix(output, inputHtml) {
		t.Fatal("Result should have prepended input html")
	}
}

func TestAppendStyleToHtmlReturnsEmptyWithNilInput(t *testing.T) {
	emptyString := ""

	css := ".test { color: red; }"
	outputPtr := AppendStyleToHtml(nil, &css)

	if outputPtr == nil {
		t.Fatal("Result should not be nil")
	}

	if (*outputPtr != emptyString) {
		t.Fatal("Result should be empty")
	}
}

func TestMergeCss(t *testing.T) {
	css1 := ".test1 { color: red; }"
	css2 := ".test2 { color: blue; }"
	css3 := ".test3 { color: green; }"

	expectedMergedCss := css1+css2+css3


	outputPtr := MergeCss(&css1, &css2, &css3)

	if outputPtr == nil {
		t.Fatal("Result should not be nil")
	}

	if *outputPtr != expectedMergedCss {
		t.Fatalf("Result do not match expectation %s != %s", *outputPtr, expectedMergedCss)
	}
}

type HttpClientMock struct {
	RequestedResources []string
}
func (c *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	c.RequestedResources = append(c.RequestedResources, req.URL.String())

    return &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}, nil
}

func TestRequestAndInlineAllHtmlResources(t *testing.T) {
	inputHtml := "<body>Hello World <img src=\"http://my-image.org/test1.png\"/> <img src=\"test2.png\"/> </body>"
	
	httpClientMock := &HttpClientMock{}
	ctx := context.WithValue(context.Background(), "httpClient", httpClientMock)

	outputPtr := RequestAndInlineAllHtmlResources(ctx, &inputHtml, "http://localhost:1234")

	if outputPtr == nil {
		t.Fatal("Result should not be nil")
	}

	if httpClientMock.RequestedResources[0] != "http://my-image.org/test1.png" {
		t.Fatal("First image was not requested")
	}

	if httpClientMock.RequestedResources[1] != "http://localhost:1234/test2.png" {
		t.Fatal("Second image was not requested")
	}

	expectedHtml := "<body>Hello World <img src=\"data:text/plain; charset=utf-8;base64,Zm9v\"/> <img src=\"data:text/plain; charset=utf-8;base64,Zm9v\"/> </body>"

	if *outputPtr != expectedHtml {
		t.Fatalf("Result was not expected: %s != %s", *outputPtr, expectedHtml)
	}
}