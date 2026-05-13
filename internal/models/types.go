package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type (
	RmgTemplate         = template.RmgTemplateModel
	RmgZone             = template.Zone
	MapTopology         = generator.MapTopology
	ZoneConfiguration   = generator.ZoneConfiguration
	AdvancedSettings    = generator.AdvancedSettings
	HeroSettings        = generator.HeroSettings
	GameEndConditions   = generator.GameEndConditions
	GladiatorArenaRules = generator.GladiatorArenaRules
	TournamentRules     = generator.TournamentRules
	GeneratorSettings   = generator.GeneratorSettings
)
