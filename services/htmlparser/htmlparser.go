package htmlparser

const (
	HeaderNodeTag = "PdfHeader"
	FooterNodeTag = "PdfFooter"
)

type HtmlParser interface {
	Parse(document *string) error
	PopHeaderAndFooter() (header string, footer string)
	AddStyles(cssStyles *string)
	GetHtml() (*string, error)
}

func New() HtmlParser {
	return &HtmlParserGoQuery{}
}
