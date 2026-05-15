package internal

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
)

func StartApplication() {
	go eventLoop()
	app.Main()
}

// eventLoop is a blocking function and needs to executed concurrently
func eventLoop() {
	window := new(app.Window)
	window.Option(
		app.Title("Olden Era — Template Generator"),
		app.Size(unit.Dp(1180), unit.Dp(820)),
		app.MinSize(unit.Dp(900), unit.Dp(600)))
	theme := themes.NewTheme()
	state := gui.NewState()

	var ops op.Ops
	for {
		switch event := window.Event().(type) {
		case app.DestroyEvent:
			if event.Err != nil {
				log.Println(event.Err)
				os.Exit(1)
			}
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)
			state.Layout(gtx, theme)
			event.Frame(gtx.Ops)
		}
	}
}
