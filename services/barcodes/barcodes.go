package barcodes

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type Barcode interface {
	GetSvg() string
}

type QrCode struct {
	barcode barcode.Barcode
}

func (qrCode *QrCode) GetSvg() string {
	return "TODO"
}

func NewQrCode(content string) (Barcode, error) {
	qr, err := qr.Encode(content, qr.M, qr.Auto)

	if err!=nil {
		return nil, err
	}

	return &QrCode{
		barcode: qr,
	}, nil
}