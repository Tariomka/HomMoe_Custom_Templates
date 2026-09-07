package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IGameRulesProvider interface {
	CreateGameRules(configuration config.GeneratorConfig) template_model.GameRules
	CreateValueOverrides(configuration config.GeneratorConfig) ([]template_model.ValueOverride, []string)
	CreateGlobalBans(configuration config.GeneratorConfig) *template_model.GlobalBans
}
