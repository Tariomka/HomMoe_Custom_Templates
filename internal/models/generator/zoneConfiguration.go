package generator

type ZoneConfiguration struct {
	NeutralZoneCount            int
	PlayerZoneCastles           int
	NeutralZoneCastles          int
	ResourceDensityPercent      int
	StructureDensityPercent     int
	NeutralStackStrengthPercent int
	BorderGuardStrengthPercent  int
	HubZoneSize                 float64
	HubZoneCastles              int
	Advanced                    AdvancedSettings
}
