package config

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type (
	AdvancedSettings    = config_inner.AdvancedSettings
	BonusEntry          = editor_state_model.BonusEntry
	BonusPresetType     = editor_state_model.BonusPresetType
	GameEndConditions   = config_inner.GameEndConditions
	GladiatorArenaRules = config_inner.GladiatorArenaRules
	HeroSettings        = config_inner.HeroSettings
	MapTopology         = topology.MapTopology
	TournamentRules     = config_inner.TournamentRules
	ZoneConfig          = config_inner.ZoneConfig
)

const (
	BonusTownPortalFree   = editor_state_model.BonusTownPortalFree
	BonusSpell            = editor_state_model.BonusSpell
	BonusUnitMultiplier   = editor_state_model.BonusUnitMultiplier
	BonusMovementBonus    = editor_state_model.BonusMovementBonus
	BonusStartingItem     = editor_state_model.BonusStartingItem
	BonusStartingGold     = editor_state_model.BonusStartingGold
	BonusStartingGems     = editor_state_model.BonusStartingGems
	BonusStartingCrystals = editor_state_model.BonusStartingCrystals
	BonusStartingMercury  = editor_state_model.BonusStartingMercury
	BonusStartingWood     = editor_state_model.BonusStartingWood
	BonusStartingOre      = editor_state_model.BonusStartingOre
)

const (
	TopologyRing         MapTopology = topology.TopologyRing
	TopologyHubAndSpoke  MapTopology = topology.TopologyHubAndSpoke
	TopologyChain        MapTopology = topology.TopologyChain
	TopologySharedWeb    MapTopology = topology.TopologySharedWeb
	TopologyRandom       MapTopology = topology.TopologyRandom
	TopologyCircles      MapTopology = topology.TopologyCircles
	TopologySquare       MapTopology = topology.TopologySquare
	TopologyGeometric    MapTopology = topology.TopologyGeometric
	TopologyCross        MapTopology = topology.TopologyCross
	TopologyFractal      MapTopology = topology.TopologyFractal
	TopologyGeometricHub MapTopology = topology.TopologyGeometricHub
)
