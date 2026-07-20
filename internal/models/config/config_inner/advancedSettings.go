package config_inner

type AdvancedSettings struct {
	NeutralLowestNoCastleCount  int
	NeutralLowestCastleCount    int
	NeutralLowestCastlesPerZone int

	NeutralLowNoCastleCount  int
	NeutralLowCastleCount    int
	NeutralLowCastlesPerZone int

	NeutralMediumNoCastleCount  int
	NeutralMediumCastleCount    int
	NeutralMediumCastlesPerZone int

	NeutralHighNoCastleCount  int
	NeutralHighCastleCount    int
	NeutralHighCastlesPerZone int

	HubZoneCastles int

	Enabled bool
}
