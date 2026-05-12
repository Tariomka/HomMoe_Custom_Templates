package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

// GeneratorSettings is the input model for the template generator. It mirrors
// the C# Olden_Era___Template_Editor.Models.GeneratorSettings type field-for-field.
//
// All values here describe a single template generation request — the GUI and
// CLI build one of these and hand it to services.Generate.
type GeneratorSettings struct {
	TemplateName string
	GameMode     string // "Classic" (only mode currently emitted; SingleHero reserved)
	PlayerCount  int    // 1..8
	MapSize      int    // raw map size in tiles (e.g. 160). sizeX = sizeZ = MapSize.

	HeroSettings *HeroSettings

	NoDirectPlayerConnections         bool // isolate player zones from each other
	RandomPortals                     bool
	MaxPortalConnections              int
	SpawnRemoteFootholds              bool
	GenerateRoads                     bool
	ExperimentalBalancedZonePlacement bool
	MatchPlayerCastleFactions         bool
	MinNeutralZonesBetweenPlayers     int

	Topology MapTopology
	ZoneCfg  ZoneConfiguration

	FactionLawsExpPercent int // 25..200 (1.0 baseline = 100)
	AstrologyExpPercent   int

	GameEndConditions   *GameEndConditions
	GladiatorArenaRules *GladiatorArenaRules
	TournamentRules     *TournamentRules

	// Optional extra mandatory content seeded by the UI; appended to the
	// player-zone defaults built by ZoneContentManager. Equivalent to the
	// C# `GeneratorSettings.PlayerZoneMandatoryContent : List<ContentItem>`.
	PlayerZoneMandatoryContent []template.MandatoryContentItem
}

// MapTopology enumerates the supported map shapes.
type MapTopology string

const (
	TopologyDefault     MapTopology = "Default" // Ring
	TopologyHubAndSpoke MapTopology = "HubAndSpoke"
	TopologyChain       MapTopology = "Chain"
	TopologySharedWeb   MapTopology = "SharedWeb"
	TopologyRandom      MapTopology = "Random"
)

// NeutralZoneQuality is the tier of a neutral zone.
type NeutralZoneQuality int

const (
	QualityLow    NeutralZoneQuality = 0
	QualityMedium NeutralZoneQuality = 1
	QualityHigh   NeutralZoneQuality = 2
)

// HeroSettings mirrors the C# HeroSettings struct.
type HeroSettings struct {
	HeroCountMin       int
	HeroCountMax       int
	HeroCountIncrement int
}

// GameEndConditions mirrors the C# GameEndConditions struct.
type GameEndConditions struct {
	VictoryCondition string // "win_condition_1"..._6
	LostStartCity    bool
	LostStartCityDay int
	LostStartHero    bool
	CityHold         bool
	CityHoldDays     int
}

// GladiatorArenaRules mirrors the C# GladiatorArenaRules struct.
type GladiatorArenaRules struct {
	Enabled        bool
	DaysDelayStart int
	CountDay       int
}

// TournamentRules mirrors the C# TournamentRules struct.
type TournamentRules struct {
	Enabled            bool
	FirstTournamentDay int
	Interval           int
	PointsToWin        int
	SaveArmy           bool
}

// ZoneConfiguration mirrors the C# ZoneConfiguration struct.
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

// AdvancedSettings mirrors the C# AdvancedSettings struct.
type AdvancedSettings struct {
	Enabled                    bool
	NeutralLowNoCastleCount    int
	NeutralLowCastleCount      int
	NeutralMediumNoCastleCount int
	NeutralMediumCastleCount   int
	NeutralHighNoCastleCount   int
	NeutralHighCastleCount     int
	PlayerZoneSize             float64
	NeutralZoneSize            float64
	GuardRandomization         float64
}

// NewGeneratorSettings returns a GeneratorSettings populated with the C# defaults.
func NewGeneratorSettings() *GeneratorSettings {
	return &GeneratorSettings{
		TemplateName:          "Custom Template",
		GameMode:              "Classic",
		PlayerCount:           2,
		MapSize:               160,
		HeroSettings:          &HeroSettings{HeroCountMin: 4, HeroCountMax: 8, HeroCountIncrement: 1},
		SpawnRemoteFootholds:  true,
		GenerateRoads:         true,
		MaxPortalConnections:  32,
		Topology:              TopologyRandom,
		FactionLawsExpPercent: 100,
		AstrologyExpPercent:   100,
		ZoneCfg: ZoneConfiguration{
			PlayerZoneCastles:           1,
			NeutralZoneCastles:          1,
			ResourceDensityPercent:      100,
			StructureDensityPercent:     100,
			NeutralStackStrengthPercent: 100,
			BorderGuardStrengthPercent:  100,
			HubZoneSize:                 1.0,
			Advanced: AdvancedSettings{
				PlayerZoneSize:     1.0,
				NeutralZoneSize:    1.0,
				GuardRandomization: 0.05,
			},
		},
		GameEndConditions:   &GameEndConditions{VictoryCondition: "win_condition_1", LostStartCityDay: 3, CityHoldDays: 6},
		GladiatorArenaRules: &GladiatorArenaRules{DaysDelayStart: 30, CountDay: 3},
		TournamentRules:     &TournamentRules{FirstTournamentDay: 14, Interval: 7, PointsToWin: 2, SaveArmy: true},
	}
}
