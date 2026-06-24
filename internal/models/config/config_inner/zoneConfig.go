package config_inner

type ZoneConfig struct {
	NeutralZoneCount            int
	PlayerOwnedCastles          int
	PlayerZoneCastles           int
	NeutralZoneCastles          int
	SpawnAbandonedOutposts      bool
	AbandonedOutpostCount       int
	ResourceDensityPercent      int
	StructureDensityPercent     int
	NeutralStackStrengthPercent int
	BorderGuardStrengthPercent  int
	HubZoneSize                 float64
	HubZoneCastles              int
	Advanced                    AdvancedSettings
}
