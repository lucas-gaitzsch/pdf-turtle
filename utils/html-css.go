package utils

import "strings"

func AppendStyleToHtml(html *string, css *string) *string {
	if html == nil {
		empty := ""
		return &empty
	}

	if css != nil {
		res := *html +
			"\n<style>" +
			*css +
			"</style>"

		return &res
	} else {
		return html
	}
}

//TODO: minify? https://github.com/tdewolff/minify

func MergeCss(css []*string) *string {
	mergedCssBuilder := strings.Builder{}

	for _, c := range css {
		mergedCssBuilder.WriteString(*c)
	}

	mergedCss := mergedCssBuilder.String()

	return &mergedCss
}
