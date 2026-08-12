package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type IGladiatorArenaProvider interface {
	PlaceArena(configuration config.GeneratorConfig, variant *entities.Variant)
}
