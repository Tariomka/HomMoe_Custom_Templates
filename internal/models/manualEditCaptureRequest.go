package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ManualEditCaptureRequest struct {
	FreshZones       []template_model.Zone
	FreshConnections []template_model.Connection
	Zones            []template_model.Zone
	Connections      []template_model.Connection
	Configuration    *config.GeneratorConfig
	NonDeterministic bool
}
