package config

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
)

type (
	AdvancedSettings    = config_inner.AdvancedSettings
	BonusEntry          = editor_state.BonusEntry
	BonusPresetType     = editor_state.BonusPresetType
	GameEndConditions   = config_inner.GameEndConditions
	GladiatorArenaRules = config_inner.GladiatorArenaRules
	HeroSettings        = config_inner.HeroSettings
	MapTopology         = editor_state.MapTopology
	TournamentRules     = config_inner.TournamentRules
	ZoneConfig          = config_inner.ZoneConfig
)

const (
	BonusTownPortalFree   = editor_state.BonusTownPortalFree
	BonusSpell            = editor_state.BonusSpell
	BonusUnitMultiplier   = editor_state.BonusUnitMultiplier
	BonusMovementBonus    = editor_state.BonusMovementBonus
	BonusStartingItem     = editor_state.BonusStartingItem
	BonusStartingGold     = editor_state.BonusStartingGold
	BonusStartingGems     = editor_state.BonusStartingGems
	BonusStartingCrystals = editor_state.BonusStartingCrystals
	BonusStartingMercury  = editor_state.BonusStartingMercury
	BonusStartingWood     = editor_state.BonusStartingWood
	BonusStartingOre      = editor_state.BonusStartingOre
)

const (
	TopologyRing         MapTopology = editor_state.TopologyRing
	TopologyHubAndSpoke  MapTopology = editor_state.TopologyHubAndSpoke
	TopologyChain        MapTopology = editor_state.TopologyChain
	TopologySharedWeb    MapTopology = editor_state.TopologySharedWeb
	TopologyRandom       MapTopology = editor_state.TopologyRandom
	TopologyCircles      MapTopology = editor_state.TopologyCircles
	TopologySquare       MapTopology = editor_state.TopologySquare
	TopologyGeometric    MapTopology = editor_state.TopologyGeometric
	TopologyCross        MapTopology = editor_state.TopologyCross
	TopologyFractal      MapTopology = editor_state.TopologyFractal
	TopologyGeometricHub MapTopology = editor_state.TopologyGeometricHub
)
