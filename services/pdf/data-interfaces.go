package pdf

type TemplateModels interface {
	GetMainModel() any
	GetHeaderModel() any
	GetFooterModel() any
	HasHeaderOrFooterModel() bool
}

type HtmlModels interface {
	HasHeaderOrFooterHtml() bool

	GetHeaderHtml() string
	SetHeaderHtml(html string)

	GetFooterHtml() string
	SetFooterHtml(html string)
}
