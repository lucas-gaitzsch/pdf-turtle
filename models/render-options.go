package models

import (
	"strings"

	"github.com/lucas-gaitzsch/pdf-turtle/utils"
)

type RenderOptionsMargins struct {
	// margin top in mm
	Top int `json:"top,omitempty" default:"25"`
	// margin right in mm
	Right int `json:"right,omitempty" default:"25"`
	// margin bottom in mm
	Bottom int `json:"bottom,omitempty" default:"20"`
	// margin left in mm
	Left int `json:"left,omitempty" default:"25"`
} // @name RenderOptionsMargins

type RenderOptions struct {
	Landscape            bool `json:"landscape,omitempty" default:"false"`
	ExcludeBuiltinStyles bool `json:"excludeBuiltinStyles,omitempty" default:"false"`

	// page size in mm; overrides page format
	PageSize   PageSize `json:"pageSize,omitempty"`
	PageFormat string   `json:"pageFormat,omitempty" default:"A4" enums:"A0,A1,A2,A3,A4,A5,A6,Letter,Legal"`

	// margins in mm; fallback to default if null
	Margins *RenderOptionsMargins `json:"margins,omitempty"`

	// true if options was parsed from bundle
	IsBundle bool `json:"-"`
	// base path is required for accessing bundle assets from loopback
	BasePath string `json:"-"`
} // @name RenderOptions

func (ro *RenderOptions) SetDefaults() {
	utils.ReflectDefaultValues(ro)

	ro.setDefaultMargin()
	ro.setEmptyPageSizeByFormat()
}

func (ro *RenderOptions) setDefaultMargin() {
	if ro.Margins == nil {
		ro.Margins = &RenderOptionsMargins{}
		utils.ReflectDefaultValues(ro.Margins)
	}
}

func (ro *RenderOptions) setEmptyPageSizeByFormat() {
	if ro.PageSize.Width == 0 || ro.PageSize.Height == 0 {
		if size, ok := PageSizesMap[strings.ToLower(ro.PageFormat)]; ok {
			ro.PageSize = size
		}
	}
}
