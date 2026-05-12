package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// labeledRowW wraps labeledRow so it can be used as a layout.Widget value
// in slice literals. Accepts a unitless dp value for the label width.
func labeledRowW(th *material.Theme, label string, labelDp int, control layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return labeledRow(gtx, th, label, unit.Dp(labelDp), control)
	}
}

// dimLabelW returns dimLabel as a layout.Widget.
func dimLabelW(th *material.Theme, txt string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return dimLabel(gtx, th, txt) }
}

// warnBannerW returns warnBanner as a layout.Widget.
func warnBannerW(th *material.Theme, txt string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return warnBanner(gtx, th, txt) }
}

// sectionHeaderW returns sectionHeader as a layout.Widget.
func sectionHeaderW(th *material.Theme, txt string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return sectionHeader(gtx, th, txt) }
}

// mapRange linearly maps a [0,1] slider value into [lo, hi].
func mapRange(v, lo, hi float32) float32 {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return lo + v*(hi-lo)
}

// mapRangeInv is the inverse of mapRange: maps a value in [lo, hi]
// back to its [0,1] slider position.
func mapRangeInv(v, lo, hi float32) float32 {
	if hi == lo {
		return 0
	}
	r := (v - lo) / (hi - lo)
	if r < 0 {
		r = 0
	}
	if r > 1 {
		r = 1
	}
	return r
}

// indexOf returns the index of v in items, or -1 when not present.
func indexOf[T comparable](items []T, v T) int {
	for i, x := range items {
		if x == v {
			return i
		}
	}
	return -1
}
