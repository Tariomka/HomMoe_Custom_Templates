package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type IZoneLayoutProvider interface {
	CreateZoneLayouts() []entities.ZoneLayoutDef
}
