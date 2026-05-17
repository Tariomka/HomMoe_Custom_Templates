package components

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type IPanel interface {
	GetWidget(theme *material.Theme) layout.Widget
}
