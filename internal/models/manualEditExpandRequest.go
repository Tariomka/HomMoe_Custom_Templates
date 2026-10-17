package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ManualEditExpandRequest struct {
	FreshZones       []template_model.Zone
	FreshConnections []template_model.Connection
	ManualEdits      editor_state_model.ManualEditSettings
	CastleChanges    editor_state_model.CastleSettingChanges
	Configuration    *config.GeneratorConfig
	NonDeterministic bool
}
