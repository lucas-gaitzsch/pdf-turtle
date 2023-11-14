package utils

import (
	"context"
	"encoding/base64"
	"io"
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

func MergeCss(css ...*string) *string {
	mergedCssBuilder := strings.Builder{}

	for _, c := range css {
		mergedCssBuilder.WriteString(*c)
	}

	mergedCss := mergedCssBuilder.String()

	return &mergedCss
}

var urlReferenceRegex = regexp.MustCompile(`( src="| href="|src: *url\(")([^"]+)("\)|")`)

type HttpClientExecuter interface {
    Do(req *http.Request) (*http.Response, error)
}

func RequestAndInlineAllHtmlResources(ctx context.Context, htmlPtr *string, baseUrl string) *string {
	if baseUrl != "" && !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	logger := zerolog.Ctx(ctx)

	html := urlReferenceRegex.ReplaceAllStringFunc(*htmlPtr, func(htmlAttribute string) string {
		return requestAndReturnBase64IfPossible(ctx, htmlAttribute, baseUrl, logger)
	})

	return &html
}

func requestAndReturnBase64IfPossible(ctx context.Context, htmlAttribute string, baseUrl string, logger *zerolog.Logger) string {

	groups := urlReferenceRegex.FindAllStringSubmatch(htmlAttribute, 2)
	prefix := groups[0][1]
	src := groups[0][2]
	suffix := groups[0][3]

	if baseUrl != "" && !strings.HasPrefix(src, "http") {
		src = baseUrl + src
	}

	var client HttpClientExecuter = http.DefaultClient
	if c, ok := ctx.Value("httpClient").(HttpClientExecuter); ok {
		client = c
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		logger.Info().Err(err).Msg("cant initialize request")
		return htmlAttribute
	}
	response, err := client.Do(req)

	if err != nil {
		logger.Info().Str("resourceUrl", src).Err(err).Msg("cant fetch resource")
		return htmlAttribute
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Info().Str("resourceUrl", src).Err(err).Msg("cant fetch resource: cant read from response body")
		return htmlAttribute
	}

	logger.Debug().Str("resourceUrl", src).Msg("resource fetched successfully")

	mimeType := http.DetectContentType(bytes)
	base64 := base64.StdEncoding.EncodeToString(bytes)

	isCssSrcUrl := strings.HasSuffix(prefix, `url("`)
	
	if (isCssSrcUrl) {
		prefix = strings.TrimRight(prefix,`"`)
		suffix = strings.TrimLeft(suffix,`"`)
	}

	return prefix + "data:" + mimeType + ";base64," + base64 + suffix
}
