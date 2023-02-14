package utils

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
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

func RequestAndInlineAllHtmlResources(ctx context.Context, htmlPtr *string, baseUrl string) *string {
	if baseUrl != "" && !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	logger := zerolog.Ctx(ctx)

	html := urlReferenceRegex.ReplaceAllStringFunc(*htmlPtr, func(htmlAttribute string) string {
		return requestAndReturnBase64IfPossible(htmlAttribute, baseUrl, logger)
	})

	return &html
}

func requestAndReturnBase64IfPossible(htmlAttribute string, baseUrl string, logger *zerolog.Logger) string {

	groups := urlReferenceRegex.FindAllStringSubmatch(htmlAttribute, 2)
	attribute := groups[0][1]
	src := groups[0][2]

	if baseUrl != "" && !strings.HasPrefix(src, "http") {
		src = baseUrl + src
	}

	response, err := http.Get(src)
	if err != nil {
		logger.Info().Str("resourceUrl", src).Err(err).Msg("cant fetch resource")
		return htmlAttribute
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Info().Str("resourceUrl", src).Err(err).Msg("cant fetch resource: cant read from response body")
		return htmlAttribute
	}

	logger.Debug().Str("resourceUrl", src).Msg("resource fetched successfully")

	mimeType := http.DetectContentType(bytes)
	base64 := base64.StdEncoding.EncodeToString(bytes)

	return attribute + "=\"data:" + mimeType + ";base64," + base64 + "\""
}
