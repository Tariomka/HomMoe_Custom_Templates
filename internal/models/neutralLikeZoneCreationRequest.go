package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"

type NeutralLikeZoneCreationRequest struct {
	Name                 string
	Profile              neutral_zone.Profile
	Size                 float64
	ConnectionNames      []string
	MandatoryContentName string
	CastleStrategy       ZoneCastleStrategy
	CastleCount          int
	HoldCity             bool
	OutpostCount         int
	FootholdCount        int
	GuardRandomization   float64
	GenerateRoads        bool
	BiomeMatchPolicy     ZoneBiomeMatchPolicy
	Tuning               GenerationTuning
}
