package models

// GeneratorSettings contains all user configuration for template generation
type GeneratorSettings struct {
	TemplateName                string
	GameMode                    string // e.g., "Classic", "Blitz", "Heroic"
	PlayerCount                 int
	MapSize                     string // "S", "M", "L", "XL", "2XL"
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
}

// HeroSettings defines hero-specific rules
type HeroSettings struct {
	HeroCount        int
	AllowedHeroes    []string
	AllowedFactions  []string
	EnableMulticlass bool
}

// TournamentRules defines tournament-specific rules
type TournamentRules struct {
	Enabled           bool
	Arenas            int
	MinimumSpacing    int
	WinningPrizeValue int
}

// GladiatorArenaRules defines gladiator arena-specific rules
type GladiatorArenaRules struct {
	Enabled      bool
	ArenasCount  int
	MinCreatures int
	MaxCreatures int
}
