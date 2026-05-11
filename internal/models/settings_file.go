package models

// ZoneContentItemUI is the editable zone-content row used by the GUI (and
// persisted in SettingsFile.PlayerZoneMandatoryContent). Mirrors
// Models/Generator/ZoneContentItemUI.cs.
type ZoneContentItemUI struct {
	Sid          string `json:"sid,omitempty"`
	Name         string `json:"name,omitempty"`
	Count        int    `json:"count"`
	IsGuarded    bool   `json:"isGuarded"`
	NearCastle   bool   `json:"nearCastle,omitempty"`
	RoadDistance string `json:"roadDistance,omitempty"`
	IsGroup      bool   `json:"isGroup,omitempty"`
}

// SettingsFile is the serialised .oetgs file produced and consumed by the
// editor. Mirrors Models/Generator/SettingsFile.cs (JSON tags identical).
type SettingsFile struct {
	TemplateName                      string              `json:"templateName"`
	MapSize                           int                 `json:"mapSize"`
	PlayerCount                       int                 `json:"playerCount"`
	NeutralZoneCount                  int                 `json:"neutralZoneCount"`
	PlayerZoneCastles                 int                 `json:"playerCastles"`
	NeutralZoneCastles                int                 `json:"neutralCastles"`
	AdvancedMode                      bool                `json:"advancedMode"`
	NeutralLowNoCastleCount           int                 `json:"neutralLowNoCastle"`
	NeutralLowCastleCount             int                 `json:"neutralLowCastle"`
	NeutralMediumNoCastleCount        int                 `json:"neutralMediumNoCastle"`
	NeutralMediumCastleCount          int                 `json:"neutralMediumCastle"`
	NeutralHighNoCastleCount          int                 `json:"neutralHighNoCastle"`
	NeutralHighCastleCount            int                 `json:"neutralHighCastle"`
	MatchPlayerCastleFactions         bool                `json:"matchPlayerCastleFactions"`
	MinNeutralZonesBetweenPlayers     int                 `json:"minNeutralZonesBetweenPlayers"`
	ExperimentalBalancedZonePlacement bool                `json:"experimentalBalancedZonePlacement"`
	ExperimentalMapSizes              bool                `json:"experimentalMapSizes"`
	PlayerZoneSize                    float64             `json:"playerZoneSize"`
	NeutralZoneSize                   float64             `json:"neutralZoneSize"`
	HubZoneSize                       float64             `json:"hubZoneSize"`
	HubZoneCastles                    int                 `json:"hubCastles"`
	GuardRandomization                float64             `json:"guardRandomization"`
	HeroCountMin                      int                 `json:"heroMin"`
	HeroCountMax                      int                 `json:"heroMax"`
	HeroCountIncrement                int                 `json:"heroIncrement"`
	Topology                          MapTopology         `json:"topology"`
	RandomPortals                     bool                `json:"randomPortals"`
	MaxPortalConnections              int                 `json:"maxPortalConns"`
	SpawnRemoteFootholds              bool                `json:"spawnFootholds"`
	GenerateRoads                     bool                `json:"generateRoads"`
	NoDirectPlayerConn                bool                `json:"isolateplayers"`
	ResourceDensityPercent            *int                `json:"resourceDensity,omitempty"`
	StructureDensityPercent           *int                `json:"structureDensity,omitempty"`
	NeutralStackStrengthPercent       int                 `json:"neutralStackStrength"`
	BorderGuardStrengthPercent        int                 `json:"borderGuardStrength"`
	VictoryCondition                  string              `json:"victoryCondition"`
	FactionLawsExpPercent             int                 `json:"factionLawsExp"`
	AstrologyExpPercent               int                 `json:"astrologyExp"`
	LostStartCity                     bool                `json:"lostStartCity"`
	LostStartCityDay                  int                 `json:"lostStartCityDay"`
	LostStartHero                     bool                `json:"lostStartHero"`
	CityHold                          bool                `json:"cityHold"`
	CityHoldDays                      int                 `json:"cityHoldDays"`
	GladiatorArena                    bool                `json:"gladiatorArena"`
	GladiatorArenaDaysDelayStart      int                 `json:"gladiatorArenaDaysDelayStart"`
	GladiatorArenaCountDay            int                 `json:"gladiatorArenaCountDay"`
	Tournament                        bool                `json:"tournament"`
	TournamentFirstTournamentDay      int                 `json:"tournamentFirstTournamentDay"`
	TournamentInterval                int                 `json:"tournamentInterval"`
	TournamentPointsToWin             int                 `json:"tournamentPointsToWin"`
	TournamentSaveArmy                bool                `json:"tournamentSaveArmy"`
	PlayerZoneMandatoryContent        []ZoneContentItemUI `json:"playerZoneMandatoryContent,omitempty"`
	ContentDensityPercent             *int                `json:"contentDensity,omitempty"`
}

// EffectiveResourceDensity matches SettingsFile.EffectiveResourceDensityPercent.
func (s *SettingsFile) EffectiveResourceDensity() int {
	if s.ResourceDensityPercent != nil {
		return *s.ResourceDensityPercent
	}
	if s.ContentDensityPercent != nil {
		return *s.ContentDensityPercent
	}
	return 100
}

// EffectiveStructureDensity matches SettingsFile.EffectiveStructureDensityPercent.
func (s *SettingsFile) EffectiveStructureDensity() int {
	if s.StructureDensityPercent != nil {
		return *s.StructureDensityPercent
	}
	if s.ContentDensityPercent != nil {
		return *s.ContentDensityPercent
	}
	return 100
}

// NewSettingsFile returns a SettingsFile populated with the C#
// SettingsFile defaults.
func NewSettingsFile() *SettingsFile {
	return &SettingsFile{
		TemplateName:                 "Custom Template",
		MapSize:                      160,
		PlayerCount:                  2,
		PlayerZoneCastles:            1,
		NeutralZoneCastles:           1,
		PlayerZoneSize:               1.0,
		NeutralZoneSize:              1.0,
		HubZoneSize:                  1.0,
		GuardRandomization:           0.05,
		HeroCountMin:                 4,
		HeroCountMax:                 8,
		HeroCountIncrement:           1,
		Topology:                     TopologyRandom,
		MaxPortalConnections:         32,
		SpawnRemoteFootholds:         true,
		GenerateRoads:                true,
		NeutralStackStrengthPercent:  100,
		BorderGuardStrengthPercent:   100,
		VictoryCondition:             "win_condition_1",
		FactionLawsExpPercent:        100,
		AstrologyExpPercent:          100,
		LostStartCityDay:             3,
		CityHoldDays:                 6,
		GladiatorArenaDaysDelayStart: 30,
		GladiatorArenaCountDay:       3,
		TournamentFirstTournamentDay: 14,
		TournamentInterval:           7,
		TournamentPointsToWin:        2,
		TournamentSaveArmy:           true,
	}
}
