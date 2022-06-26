package utils

import (
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
