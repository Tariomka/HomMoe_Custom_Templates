package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Dialog is a modal view rendered on top of the main window by the DialogHost.
// Implementations draw only their body; the host supplies the scrim, the
// centered panel chrome, and a title bar with a close button.
//
// The interface lives in the leaf widgets package so that deep components (e.g.
// the zone-content rows) can construct dialogs without importing the parent
// components package, which would create an import cycle.
type Dialog interface {
	// Title is shown in the modal header.
	Title() string
	// Body draws the dialog content inside the panel and returns true when the
	// dialog has finished and should be closed.
	Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool)
	// PreferredSize is the desired panel size in Dp. A zero component falls back
	// to the host default.
	PreferredSize() (width, height unit.Dp)
}

// DialogOpener opens a modal dialog. The DialogHost.Open method value satisfies
// this type, letting components request dialogs without a direct dependency on
// the host's concrete type.
type DialogOpener = func(Dialog)
