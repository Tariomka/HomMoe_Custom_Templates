package provider_interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"

type IZoneLayoutProvider interface {
	CreateZoneLayouts() []template_model.ZoneLayoutDef
}
