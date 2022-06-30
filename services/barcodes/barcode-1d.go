package barcodes

import (
	"image/color"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode"
)

const aspectRatioBarcode1d = 10 / 4

type Barcode1D struct {
	data barcode.Barcode
}

func (bc *Barcode1D) Svg() string {
	bcd := bc.data

	bounds := bcd.Bounds()

	qrWidth := bounds.Dx()

	svgWidth := qrWidth
	svgHeight := qrWidth / aspectRatioBarcode1d

	sb := new(strings.Builder)
	svg := svg.New(sb)

	initSvg(svg, svgWidth, svgHeight)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		if bcd.At(x, bounds.Min.Y) == color.Black {
			start := x
			x++

			for x < bounds.Max.X && bcd.At(x, bounds.Min.Y) == color.Black {
				x++
			}

			if start != x {
				svg.Rect(start, 0, x-start, svgHeight, svgBlackRectangleAttr)
			}
		}
	}

	svg.End()

	return sb.String()
}
