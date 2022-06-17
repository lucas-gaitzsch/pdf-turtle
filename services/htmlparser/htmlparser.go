package htmlparser

const (
	HeaderNodeTag = "PdfHeader"
	FooterNodeTag = "PdfFooter"
	trimCutset    = "\t\n "
)

type HtmlParser interface {
	Parse(document *string) error
	PopHeaderAndFooter() (header string, footer string)
	GetHtml() (*string, error)
}

func New() HtmlParser {
	return &HtmlParserGoQuery{}
}

//TODO:! add style