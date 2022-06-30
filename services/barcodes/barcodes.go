package barcodes

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
)

type BarcodeSvg interface {
	Svg() string
}

func NewBarcodeSvg(barcodeCreationFunc func() (barcode.Barcode, error)) (BarcodeSvg, error) {
	bc, err := barcodeCreationFunc()

	if err != nil {
		return nil, err
	}

	if bc.Metadata().Dimensions == 1 {
		return &Barcode1D{
			data: bc,
		}, nil
	} else {
		return &Barcode2D{
			data: bc,
		}, nil
	}
}

func NewEanCode(content string) (BarcodeSvg, error) {
	return NewBarcodeSvg(func() (barcode.Barcode, error) { return ean.Encode(content) })
}

func NewQrCode(content string) (BarcodeSvg, error) {
	return NewBarcodeSvg(func() (barcode.Barcode, error) { return qr.Encode(content, qr.L, qr.Auto) })
}
