package gui

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/editor"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
)

func StartApplication(version string) {
	go eventLoop(version)
	app.Main()
}

// eventLoop is a blocking function and needs to executed concurrently.
func eventLoop(version string) {
	window := getAndConfigureWindow(version)
	windowLayout := editor.NewWindow(
		composition.InitializeGuiHandler(),
		composition.InitializeFileSystemHandler(),
		composition.InitializeRegenerationHandler())
	windowLayout.SetOnExit(func() { window.Perform(system.ActionClose) })
	theme := themes.NewTheme()
	positionLogger := utils.NewButtonPositionLogger(slog.Default())

	var ops op.Ops
	for {
		switch event := window.Event().(type) {
		case app.DestroyEvent:
			if event.Err != nil {
				slog.Error("Window destroyed with error", slog.String("error", event.Err.Error()))
				os.Exit(1)
			}
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)
			windowLayout.Layout(gtx, theme)
			event.Frame(gtx.Ops)
			positionLogger.LogButtonPositions(gtx.Ops, gtx.Metric, event.Size)
		}
	}
}

func getAndConfigureWindow(version string) *app.Window {
	window := new(app.Window)
	configuration := []app.Option{
		app.Title("Olden Era - Custom Templates (" + version + ")"),
		app.MinSize(unit.Dp(1280), unit.Dp(800)),
	}

	slog.SetDefault(slog.New(slog.DiscardHandler))

	windowWidth, windowHeight := 1600, 900
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-minimized":
			configuration = append(configuration, app.Minimized.Option())
			continue
		case "-fullscreen":
			configuration = append(configuration, app.Fullscreen.Option())
			continue
		case "-with-logging":
			handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
			slog.SetDefault(slog.New(handler))
			continue
		}

		split := strings.Split(arg, "=")
		if len(split) > 1 {
			if split[1] == "" {
				slog.Warn("Value for argument is missing", slog.String("argument", split[0]))
			}

			switch split[0] {
			case "-w":
				if width, err := strconv.Atoi(split[1]); err == nil {
					windowWidth = width
				}
				continue
			case "-h":
				if height, err := strconv.Atoi(split[1]); err == nil {
					windowHeight = height
				}
				continue
			}
		}
	}

	configuration = append(
		// Order matters: Minimized / Fullscreen get overridden by Size
		[]app.Option{app.Size(unit.Dp(windowWidth), unit.Dp(windowHeight))},
		configuration...,
	)
	window.Option(configuration...)
	return window
}
