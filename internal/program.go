package application

import (
	"log"
	"os"

	"gioui.org/app"

	"github.com/Tariomka/hommoe_custom_templates/internal/gui"
)

func Start() {
	go func() {
		if err := gui.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}
