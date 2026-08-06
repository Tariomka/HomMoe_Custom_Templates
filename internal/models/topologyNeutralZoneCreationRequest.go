package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"

type TopologyNeutralZoneCreationRequest struct {
	Plan            neutral_zone.Plan
	ConnectionNames []string
	Size            float64
	FootholdCount   int
	GenerateRoads   bool
	HoldCity        bool
	Tuning          GenerationTuning
}
