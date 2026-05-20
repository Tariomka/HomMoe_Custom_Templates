package generator

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

// GeneratorSettings is the input model for the template generator
//
// All values here describe a single template generation request — the GUI and
// CLI build one of these and hand it to services.Generate.
type GeneratorSettings struct {
	TemplateName string
	GameMode     string // "Classic" (only mode currently emitted; SingleHero reserved)
	PlayerCount  int    // 1..8
	MapSize      int    // raw map size in tiles (e.g. 160). sizeX = sizeZ = MapSize.

	HeroSettings *HeroSettings

	NoDirectPlayerConnections     bool // isolate player zones from each other
	RandomPortals                 bool
	MaxPortalConnections          int
	SpawnRemoteFootholds          bool
	GenerateRoads                 bool
	MatchPlayerCastleFactions     bool
	MinNeutralZonesBetweenPlayers int

	// Banned content & overrides (raw strings as edited in the UI).
	BannedItems        string
	BannedMagics       string
	ValueOverridesText string

	// Configurable game-start bonuses (Wood/Ore/spell/etc.). Parsed from
	// SettingsFile.BonusesJson by the loader.
	Bonuses []BonusEntry

	Topology MapTopology
	ZoneCfg  ZoneConfiguration

	FactionLawsExpPercent int // 25..200 (1.0 baseline = 100)
	AstrologyExpPercent   int

	GameEndConditions   *GameEndConditions
	GladiatorArenaRules *GladiatorArenaRules
	TournamentRules     *TournamentRules

	// Optional extra mandatory content seeded by the UI; appended to the
	// player-zone defaults built by ZoneContentManager.
	PlayerZoneMandatoryContent    []template.MandatoryContentItem
	LowNeutralMandatoryContent    []template.MandatoryContentItem
	MediumNeutralMandatoryContent []template.MandatoryContentItem
	HighNeutralMandatoryContent   []template.MandatoryContentItem
	HubZoneMandatoryContent       []template.MandatoryContentItem
}

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
		Topology:              TopologyBalanced,
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
