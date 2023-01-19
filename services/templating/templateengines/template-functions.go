package templateengines

import (
	"encoding/json"
	"html/template"
	"strings"

	"github.com/lucas-gaitzsch/pdf-turtle/services/barcodes"
)

var templateFunctions = template.FuncMap{
	"marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
	"barcodeQr": func(content string) template.HTML {
		qr, _ := barcodes.NewQrCode(content)
		return template.HTML(qr.Svg())
	},
	"barcodeEan": func(content string) template.HTML {
		qr, _ := barcodes.NewEanCode(content)
		return template.HTML(qr.Svg())
	},
	"contains":  strings.Contains,
    "hasPrefix": strings.HasPrefix,
    "hasSuffix": strings.HasSuffix,
}
