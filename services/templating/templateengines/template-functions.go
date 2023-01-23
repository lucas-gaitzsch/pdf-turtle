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
	"strContains":  strings.Contains,
	"strHasPrefix": strings.HasPrefix,
	"strHasSuffix": strings.HasSuffix,
	"add": func(a float64, b float64) float64 {
		return a + b
	},
	"subtract": func(a float64, b float64) float64 {
		return a - b
	},
	"multiply": func(a float64, b float64) float64 {
		return a * b
	},
	"divide": func(a float64, b float64) float64 {
		return a / b
	},
	"float64ToInt": func(val float64) int {
		return int(val)
	},
	"intToFloat64": func(val int) float64 {
		return float64(val)
	},
	"bitwiseAnd": func(a int, b int) int {
		return a & b
	},
}
