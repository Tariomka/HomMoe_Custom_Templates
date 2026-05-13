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
