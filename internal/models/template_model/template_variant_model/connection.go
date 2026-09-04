package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
)

type Connection struct {
	Name string
	From string
	To   string

	ConnectionType string

	SimTurnSquad bool
	Road         *bool

	GuardZone   string
	GuardEscape bool

	GuardValue           int
	GuardRandomization   float64
	GuardWeeklyIncrement float64

	GatePlacement string

	Length float64

	GuardMatchGroup string

	PortalPlacementRulesFrom []template_common_model.PlacementRule
	PortalPlacementRulesTo   []template_common_model.PlacementRule

	// IsUserAdded marks a connection added by hand in the zone editor rather
	// than produced by the generator. It has no entity counterpart: the editor
	// state persists it on editor_state.ManualConnectionSave instead.
	IsUserAdded bool
}

func ToConnectionModel(entity template.Connection) Connection {
	return Connection{
		Name:                     entity.Name,
		From:                     entity.From,
		To:                       entity.To,
		ConnectionType:           entity.ConnectionType,
		SimTurnSquad:             entity.SimTurnSquad,
		Road:                     entity.Road,
		GuardZone:                entity.GuardZone,
		GuardEscape:              entity.GuardEscape,
		GuardValue:               entity.GuardValue,
		GuardRandomization:       entity.GuardRandomization,
		GuardWeeklyIncrement:     entity.GuardWeeklyIncrement,
		GatePlacement:            entity.GatePlacement,
		Length:                   entity.Length,
		GuardMatchGroup:          entity.GuardMatchGroup,
		PortalPlacementRulesFrom: template_common_model.ToPlacementRuleModels(entity.PortalPlacementRulesFrom),
		PortalPlacementRulesTo:   template_common_model.ToPlacementRuleModels(entity.PortalPlacementRulesTo),
	}
}

func ToConnectionEntity(model Connection) template.Connection {
	return template.Connection{
		Name:                     model.Name,
		From:                     model.From,
		To:                       model.To,
		ConnectionType:           model.ConnectionType,
		SimTurnSquad:             model.SimTurnSquad,
		Road:                     model.Road,
		GuardZone:                model.GuardZone,
		GuardEscape:              model.GuardEscape,
		GuardValue:               model.GuardValue,
		GuardRandomization:       model.GuardRandomization,
		GuardWeeklyIncrement:     model.GuardWeeklyIncrement,
		GatePlacement:            model.GatePlacement,
		Length:                   model.Length,
		GuardMatchGroup:          model.GuardMatchGroup,
		PortalPlacementRulesFrom: template_common_model.ToPlacementRuleEntities(model.PortalPlacementRulesFrom),
		PortalPlacementRulesTo:   template_common_model.ToPlacementRuleEntities(model.PortalPlacementRulesTo),
	}
}

func ToConnectionModels(entities []template.Connection) []Connection {
	return helpers.MapSlice(entities, ToConnectionModel)
}

func ToConnectionEntities(models []Connection) []template.Connection {
	return helpers.MapSlice(models, ToConnectionEntity)
}
