package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"

type RoadReconciliationRequest struct {
	Zones            []template_model.Zone
	Connections      []template_model.Connection
	MandatoryContent []template_model.MandatoryContent // Empty means the zones own no content, so every road pointing at a named content item is dropped.

	GenerateRoads        bool
	SpawnRemoteFootholds bool
	RemoteFootholdCount  int
}
