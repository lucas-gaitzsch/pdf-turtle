package models

type PageSize struct {
	// in mm
	Width int `json:"width" example:"210"`
	// in mm
	Height int `json:"height" example:"297"`
} // @name PageSize

const (
	// DIN dimensions
	PageSizeKeyA0 = "a0"
	PageSizeKeyA1 = "a1"
	PageSizeKeyA2 = "a2"
	PageSizeKeyA3 = "a3"
	PageSizeKeyA4 = "a4"
	PageSizeKeyA5 = "a5"
	PageSizeKeyA6 = "a6"

	// US dimensions
	PageSizeKeyLetter = "letter"
	PageSizeKeyLegal  = "legal"
)

var PageSizesMap = map[string]PageSize{
	// DIN dimensions
	PageSizeKeyA0: {Width: 841, Height: 1189},
	PageSizeKeyA1: {Width: 594, Height: 841},
	PageSizeKeyA2: {Width: 420, Height: 594},
	PageSizeKeyA3: {Width: 297, Height: 420},
	PageSizeKeyA4: {Width: 210, Height: 297},
	PageSizeKeyA5: {Width: 148, Height: 210},
	PageSizeKeyA6: {Width: 105, Height: 148},

	// US dimensions
	PageSizeKeyLetter: {Width: 216, Height: 297},
	PageSizeKeyLegal:  {Width: 216, Height: 279},
}
