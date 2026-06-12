package interfaces

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Dialog is a modal view rendered on top of the main window by the DialogHost.
// Implementations draw only their body; the host supplies the scrim, the
// centered panel chrome, and a title bar with a close button.
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
