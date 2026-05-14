package gui

import (
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

// Run opens the GUI window and runs the event loop until close.
func Run() error {
	window := new(app.Window)
	window.Option(
		app.Title("Olden Era — Template Generator"),
		app.Size(unit.Dp(1180), unit.Dp(820)),
		app.MinSize(unit.Dp(900), unit.Dp(600)),
	)
	theme := NewTheme()
	state := newState()

	var ops op.Ops
	for {
		switch event := window.Event().(type) {
		case app.DestroyEvent:
			return event.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)
			state.Layout(gtx, theme)
			event.Frame(gtx.Ops)
		}
	}
}
