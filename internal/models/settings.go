package models

// GeneratorSettings contains all user configuration for template generation.
type GeneratorSettings struct {
	TemplateName                string
	GameMode                    string // e.g., "Classic", "SingleHero"
	PlayerCount                 int
	MapSize                     string // "S", "M", "L", "XL", "2XL"
	MapSizeValue                int    // Raw map size, e.g. 160. Mirrors C# MapSize.
	Topology                    MapTopology
	AllowRoads                  bool
	AllowPortals                bool
	AllowFootholds              bool
	EnablePlayerIsolation       bool
	EnableCityHold              bool
	ZoneConfigurations          map[int]*ZoneConfiguration // keyed by player number
	AdvancedSettings            *AdvancedSettings
	GameEndConditions           *GameEndConditions
	HeroSettings                *HeroSettings
	TournamentRules             *TournamentRules
	GladiatorArenaRules         *GladiatorArenaRules
	ShowDescription             bool
	IncludeOptionsInDescription bool

	// Extended C# parity fields.
	MaxPortalConnections          int
	MinNeutralZonesBetweenPlayers int
	MatchPlayerCastleFactions     bool
	NeutralStackStrengthPercent   int
	BorderGuardStrengthPercent    int
	ResourceDensityPercent        int
	StructureDensityPercent       int
	FactionLawsExpPercent         int
	AstrologyExpPercent           int
	PlayerZoneSize                float64
	NeutralZoneSize               float64
	HubZoneSize                   float64
	HubZoneCastles                int
	PlayerZoneMandatoryContent    []ZoneContentItemUI
}

// MapTopology defines the topology type
type MapTopology string

const (
	TopologyDefault     MapTopology = "Default" // Ring topology
	TopologyChain       MapTopology = "Chain"
	TopologyHubAndSpoke MapTopology = "HubAndSpoke"
	TopologySharedWeb   MapTopology = "SharedWeb"
	TopologyRandom      MapTopology = "Random"
)

// ZoneConfiguration defines per-zone settings
type ZoneConfiguration struct {
	CastleCount       int     // Expected castle count
	DensityPercentage float64 // Density for content placement
}

// NeutralZoneQuality defines neutral zone tiers
type NeutralZoneQuality string

const (
	QualityLow    NeutralZoneQuality = "Low"
	QualityMedium NeutralZoneQuality = "Medium"
	QualityHigh   NeutralZoneQuality = "High"
)

// AdvancedSettings contains advanced configuration
type AdvancedSettings struct {
	NeutralZoneCount       int
	NeutralZoneLowCount    int
	NeutralZoneMediumCount int
	NeutralZoneHighCount   int
	PlayerCastleCount      int
	NeutralCastleCount     int
	GuardRandomization     float64 // Default 0.05
	ContentScaling         float64 // Default 1.0
	ConnectionCountPerZone int     // Default 2

	// Extended C# parity — split-by-castle counts.
	Enabled                    bool
	NeutralLowNoCastleCount    int
	NeutralLowCastleCount      int
	NeutralMediumNoCastleCount int
	NeutralMediumCastleCount   int
	NeutralHighNoCastleCount   int
	NeutralHighCastleCount     int
}

// GameEndConditions defines win/loss conditions
type GameEndConditions struct {
	EnableClassicVictory bool
	EnableDesertion      bool
	EnableHeroLighting   bool
	EnableCityHold       bool
	CityHoldTownType     string
	EnableGladiatorArena bool
	EnableTournaments    bool

	// Extended C# parity.
	VictoryCondition string // e.g. "win_condition_1"
	LostStartCity    bool
	LostStartCityDay int
	LostStartHero    bool
	CityHoldDays     int
}

// HeroSettings defines hero-specific rules
type HeroSettings struct {
	HeroCount        int
	AllowedHeroes    []string
	AllowedFactions  []string
	EnableMulticlass bool

	// Extended C# parity.
	HeroCountMin       int
	HeroCountMax       int
	HeroCountIncrement int
}

// TournamentRules defines tournament-specific rules
type TournamentRules struct {
	Enabled           bool
	Arenas            int
	MinimumSpacing    int
	WinningPrizeValue int

	// Extended C# parity.
	FirstTournamentDay int
	Interval           int
	PointsToWin        int
	SaveArmy           bool
}

// GladiatorArenaRules defines gladiator arena-specific rules
type GladiatorArenaRules struct {
	Enabled      bool
	ArenasCount  int
	MinCreatures int
	MaxCreatures int

	// Extended C# parity.
	DaysDelayStart int
	CountDay       int
}
