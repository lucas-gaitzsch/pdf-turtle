package barcodes

import (
	"image/color"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode"
)

type Barcode2D struct {
	data barcode.Barcode
}

func (bc *Barcode2D) Svg() string {
	bcd := bc.data

	blockSize := 10

	bounds := bcd.Bounds()

	qrWidth := bounds.Dx()

	svgWidth := qrWidth * blockSize

	sb := new(strings.Builder)
	svg := svg.New(sb)

	initSvg(svg, svgWidth, svgWidth)

	svgX := 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {

		svgY := 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			startSvgY := svgY
			for y < bounds.Max.Y && bcd.At(x, y) == color.Black {
				y++
				svgY += blockSize
			}

			if startSvgY != svgY {
				svg.Rect(svgX, startSvgY, blockSize, svgY-startSvgY, svgBlackRectangleAttr)
			}

			// if bcd.At(x, y) == color.Black {
			// 	svg.Rect(svgX, svgY, blockSize, blockSize, svgBlackRectangleAttr)
			// }
			svgY += blockSize
		}
		svgX += blockSize
	}

	svg.End()

	return sb.String()
}
