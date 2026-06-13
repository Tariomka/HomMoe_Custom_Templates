package gui

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/editor"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

func StartApplication() {
	go eventLoop()
	app.Main()
}

// eventLoop is a blocking function and needs to executed concurrently
func eventLoop() {
	window := new(app.Window)
	window.Option(
		app.Title("Olden Era - Custom Template Editor"),
		app.Size(unit.Dp(1600), unit.Dp(900)),
		app.MinSize(unit.Dp(1280), unit.Dp(800)))
	if os.Getenv("HOT_RELOAD") == "1" {
		window.Option(app.Minimized.Option())
	}

	theme := themes.NewTheme()
	windowLayout := editor.NewWindow()

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
			windowLayout.Layout(gtx, theme)
			event.Frame(gtx.Ops)
		}
	}
}
