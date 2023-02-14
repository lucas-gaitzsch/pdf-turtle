package utils

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func AppendStyleToHtml(html *string, css *string) *string {
	if html == nil {
		empty := ""
		return &empty
	}

	if css != nil {
		styleSection :=
			"\n<style>" +
				*css +
				"</style>"

		res := *html + styleSection
		return &res
	} else {
		return html
	}
}

func MergeCss(css []*string) *string {
	mergedCssBuilder := strings.Builder{}

	for _, c := range css {
		mergedCssBuilder.WriteString(*c)
	}

	mergedCss := mergedCssBuilder.String()

	return &mergedCss
}

var urlReferenceRegex = regexp.MustCompile(` (src|href)="([^"]+)"`)

func RequestAndInlineAllHtmlResources(htmlPtr *string, baseUrl string) *string {	
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	html := urlReferenceRegex.ReplaceAllStringFunc(*htmlPtr, func(htmlAttribute string)string { return requestAndReturnBase64IfPossible(htmlAttribute, baseUrl) })
	
	return &html
}

func requestAndReturnBase64IfPossible(htmlAttribute string, baseUrl string) string {
	groups := urlReferenceRegex.FindAllStringSubmatch(htmlAttribute, 2)
	attribute := groups[0][1]
	src := groups[0][2]

	if !strings.HasPrefix(src, "http") {
		src = baseUrl + src
	}
	
	response, err := http.Get(src)
	if err != nil {
		return htmlAttribute
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return htmlAttribute
	}

	mimeType := http.DetectContentType(bytes)
	base64 := base64.StdEncoding.EncodeToString(bytes)

	return attribute + "=\"data:" + mimeType + ";base64," + base64 + "\""
}
