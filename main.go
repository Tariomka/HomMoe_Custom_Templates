package main

import (
	application "github.com/Tariomka/hommoe_custom_templates/app/gui"
)

// version is injected at release build time via -ldflags "-X main.version=vX.Y.Z".
var version = "dev"

func main() {
	application.StartApplication(version)
}
