package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"

type NeutralZoneCreationRequest struct {
	Name                 string
	Quality              neutral_zone.Quality
	Size                 float64
	ConnectionNames      []string
	MandatoryContentName string
	CastleCount          int
	HoldCity             bool
	OutpostCount         int
	FootholdCount        int
	GuardRandomization   float64
	GenerateRoads        bool
	Tuning               GenerationTuning
}
