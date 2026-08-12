package models

type SpawnZoneCreationRequest struct {
	Label           string
	PlayerName      string
	ConnectionNames []string
	CastleCount     int
	MatchFactions   bool
	Size            float64
	FootholdCount   int
	GenerateRoads   bool
	Tuning          GenerationTuning
}
