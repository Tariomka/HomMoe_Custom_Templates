// Command gui launches the Olden Era Template Generator desktop UI.
//
// The UI is rendered with Gio (https://gioui.org). On Windows, Gio uses the
// Direct3D 11 backend through syscall — pure Go, no CGO, no OpenGL required.
package main

import (
	"log"
	"os"

	"gioui.org/app"

	gui "github.com/Tariomka/hommoe_custom_templates/internal/gui"
)

func main() {
	go func() {
		if err := gui.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}
