package config

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

const (
	defaultTemplateName = "Custom Template"
)

var (
	winConditions = registry.GetWinningConditionValues()
	gameModes     = registry.GetGameModeValues()
)

// GeneratorConfig is the input model for the template generator
//
// All values here describe a single template generation request - the GUI and
// CLI build one of these and hand it to services.Generate.
type GeneratorConfig struct {
	TemplateName string
	GameMode     string // "Classic"/"SingleHero"
	PlayerCount  int    // 1..8
	MapSize      int    // raw map size in tiles (e.g. 160). sizeX = sizeZ = MapSize.

	HeroSettings HeroSettings

	NoDirectPlayerConnections     bool // isolate player zones from each other
	RandomPortals                 bool
	MaxPortalConnections          int
	SpawnRemoteFootholds          bool
	RemoteFootholdCount           int
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

	Topology          MapTopology
	ZoneConfiguration ZoneConfig

	FactionLawsExpPercent int // 25..200 (1.0 baseline = 100)
	AstrologyExpPercent   int

	GameEndConditions   *GameEndConditions
	GladiatorArenaRules *GladiatorArenaRules
	TournamentRules     *TournamentRules

	// Optional extra mandatory content seeded by the UI; appended to the
	// player-zone defaults built by ZoneContentManager.
	PlayerZoneMandatoryContent    []entities.MandatoryContentItem
	LowNeutralMandatoryContent    []entities.MandatoryContentItem
	MediumNeutralMandatoryContent []entities.MandatoryContentItem
	HighNeutralMandatoryContent   []entities.MandatoryContentItem
	HubZoneMandatoryContent       []entities.MandatoryContentItem

	// ShufflePlayerZones randomizes which physical zone each player starts in.
	// Enabled by default so generated templates vary between runs; tests can
	// disable it to obtain deterministic output.
	ShufflePlayerZones bool
}

func NewGeneratorConfig() *GeneratorConfig {
	return &GeneratorConfig{
		TemplateName: defaultTemplateName,
		GameMode:     gameModes.Classic,
		PlayerCount:  2,
		MapSize:      160,
		HeroSettings: HeroSettings{
			HeroCountMin:       4,
			HeroCountMax:       8,
			HeroCountIncrement: 1,
		},
		SpawnRemoteFootholds:  true,
		RemoteFootholdCount:   1,
		GenerateRoads:         true,
		MaxPortalConnections:  32,
		Topology:              config_inner.TopologyCircles,
		FactionLawsExpPercent: 100,
		AstrologyExpPercent:   100,
		ZoneConfiguration: ZoneConfig{
			PlayerZoneCastles:           1,
			NeutralZoneCastles:          1,
			AbandonedOutpostCount:       1,
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
		GameEndConditions: &GameEndConditions{
			VictoryCondition: winConditions.Standard,
			LostStartCityDay: 3,
			CityHoldDays:     6,
		},
		GladiatorArenaRules: &GladiatorArenaRules{DaysDelayStart: 30, CountDay: 3},
		TournamentRules: &TournamentRules{
			FirstTournamentDay: 14,
			Interval:           7,
			PointsToWin:        2,
			SaveArmy:           true,
		},
		ShufflePlayerZones: true,
	}
}

func (this *GeneratorConfig) IsTournamentMode() bool {
	return (this.TournamentRules != nil && this.TournamentRules.Enabled) ||
		(this.GameEndConditions != nil && this.GameEndConditions.VictoryCondition == winConditions.Tournament)
}

func (this *GeneratorConfig) IsCityHoldMode() bool {
	return this.GameEndConditions != nil &&
		(this.GameEndConditions.CityHold || this.GameEndConditions.VictoryCondition == winConditions.CityHold)
}

func (this *GeneratorConfig) IsHubCityToHold() bool {
	return this.Topology == config_inner.TopologyHubAndSpoke && this.IsCityHoldMode()
}

func (this *GeneratorConfig) IsSingleHeroMode() bool {
	return this.GameMode == gameModes.SingleHero
}

func (this *GeneratorConfig) GetVictoryCondition() string {
	if this.GameEndConditions != nil {
		return this.GameEndConditions.VictoryCondition
	}

	return winConditions.Standard
}

func (this *GeneratorConfig) GetHeroSettings() HeroSettings {
	if this.IsSingleHeroMode() {
		return HeroSettings{
			HeroCountMin:       1,
			HeroCountMax:       1,
			HeroCountIncrement: 1,
		}
	}
	return this.HeroSettings
}

func (this *GeneratorConfig) GetGameEndConditions() GameEndConditions {
	if this.GameEndConditions != nil {
		return *this.GameEndConditions
	}
	return GameEndConditions{
		VictoryCondition: winConditions.Standard,
		LostStartCityDay: 3,
		CityHoldDays:     6,
	}
}

func (this *GeneratorConfig) GetGladiatorArenaRules() GladiatorArenaRules {
	if this.GladiatorArenaRules != nil {
		return *this.GladiatorArenaRules
	}
	return GladiatorArenaRules{}
}

func (this *GeneratorConfig) GetTournamentRules() TournamentRules {
	if this.TournamentRules != nil {
		return *this.TournamentRules
	}
	return TournamentRules{}
}

func (this *GeneratorConfig) EnsureNameExists() {
	if this.TemplateName == "" {
		this.TemplateName = defaultTemplateName
	}
}

func (this *GeneratorConfig) CanHonorNeutralSeparation() bool {
	min := this.MinNeutralZonesBetweenPlayers
	if min <= 0 {
		return true
	}

	if this.RandomPortals {
		return false
	}

	neutralZoneCount := this.getNeutralZoneCount()
	switch this.Topology {
	case config_inner.TopologyRing, config_inner.TopologyCircles:
		return neutralZoneCount >= this.PlayerCount*min
	case config_inner.TopologyChain:
		return neutralZoneCount >= (this.PlayerCount-1)*min
	case config_inner.TopologyHubAndSpoke:
		return min <= 1
	case config_inner.TopologySharedWeb:
		return min <= 1 && neutralZoneCount >= 1
	default:
		return false
	}
}

func (this *GeneratorConfig) getNeutralZoneCount() int {
	advancedTotal := this.ZoneConfiguration.Advanced.NeutralLowNoCastleCount +
		this.ZoneConfiguration.Advanced.NeutralLowCastleCount +
		this.ZoneConfiguration.Advanced.NeutralMediumNoCastleCount +
		this.ZoneConfiguration.Advanced.NeutralMediumCastleCount +
		this.ZoneConfiguration.Advanced.NeutralHighNoCastleCount +
		this.ZoneConfiguration.Advanced.NeutralHighCastleCount

	count := helpers.BoolToInt(this.Topology == config_inner.TopologySharedWeb)
	if advancedTotal > 0 {
		count += advancedTotal
	} else {
		count += this.ZoneConfiguration.NeutralZoneCount
	}

	return count
}
