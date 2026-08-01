package models

type HubZoneCreation struct {
	Name                 string
	Size                 float64
	ConnectionNames      []string
	MandatoryContentName string
	CastleCount          int
	HoldCity             bool
	GuardRandomization   float64
	GenerateRoads        bool
	Tuning               GenerationTuning
}
