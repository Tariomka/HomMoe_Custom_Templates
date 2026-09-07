package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IGladiatorArenaProvider interface {
	PlaceArena(configuration config.GeneratorConfig, variant *template_model.Variant)
}
