package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// labeledRowW wraps labeledRow so it can be used as a layout.Widget value
// in slice literals. Accepts a unitless dp value for the label width.
func labeledRowW(theme *material.Theme, label string, labelDp int, control layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return labeledRow(gtx, theme, label, unit.Dp(labelDp), control)
	}
}

// dimLabelW returns dimLabel as a layout.Widget.
func dimLabelW(theme *material.Theme, txt string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return dimLabel(gtx, theme, txt) }
}

// warnBannerW returns warnBanner as a layout.Widget.
func warnBannerW(theme *material.Theme, txt string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return warnBanner(gtx, theme, txt) }
}

// sectionHeaderW returns sectionHeader as a layout.Widget.
func sectionHeaderW(theme *material.Theme, txt string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return sectionHeader(gtx, theme, txt) }
}

// mapRange linearly maps a [0,1] slider value into [low, high].
func mapRange(value, low, high float32) float32 {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return low + value*(high-low)
}

// mapRangeInv is the inverse of mapRange: maps a value in [low, high]
// back to its [0,1] slider position.
func mapRangeInv(value, low, high float32) float32 {
	if high == low {
		return 0
	}
	ratio := (value - low) / (high - low)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return ratio
}

// indexOf returns the index of value in items, or -1 when not present.
func indexOf[T comparable](items []T, value T) int {
	for i, candidate := range items {
		if candidate == value {
			return i
		}
	}
	return -1
}

// mapSizeLabelInt returns the short S/M/L/XL/H/G/C label for an integer size.
// Mirrors KnownValues.MapSizeLabel from the C# source.
func mapSizeLabelInt(size int) string {
	switch {
	case size == 64:
		return "S"
	case size == 80 || size == 96:
		return "M"
	case size == 112 || size == 128:
		return "L"
	case size == 144 || size == 160:
		return "XL"
	case size == 176 || size == 192:
		return "H"
	case size >= 208 && size <= 256:
		return "G"
	default:
		return "C"
	}
}
