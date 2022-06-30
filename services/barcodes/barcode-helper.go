package barcodes

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
)

const svgBlackRectangleAttr = "fill:black; stroke:none;"

func initSvg(svg *svg.SVG, width int, height int) *svg.SVG {
	svg.Start(
		width,
		height,
		fmt.Sprintf(`viewBox="0 0 %d %d"`, width, height),
		`style="width:100%; height: 100%;"`,
	)
	return svg
}
