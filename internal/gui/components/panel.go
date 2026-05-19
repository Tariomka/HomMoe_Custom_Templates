package components

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type IPanel interface {
	GetPanelWidget(theme *material.Theme) layout.Widget
	LoadFromState()
	SaveToState()
}
