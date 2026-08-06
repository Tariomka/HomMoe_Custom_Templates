package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type IGameRulesProvider interface {
	CreateGameRules(configuration config.GeneratorConfig) entities.GameRules
	CreateValueOverrides(configuration config.GeneratorConfig) ([]entities.ValueOverride, []string)
	CreateGlobalBans(configuration config.GeneratorConfig) *entities.GlobalBans
}
