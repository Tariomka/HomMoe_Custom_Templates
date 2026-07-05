package config

import "github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"

type (
	AdvancedSettings    = config_inner.AdvancedSettings
	BonusEntry          = config_inner.BonusEntry
	BonusPresetType     = config_inner.BonusPresetType
	GameEndConditions   = config_inner.GameEndConditions
	GladiatorArenaRules = config_inner.GladiatorArenaRules
	HeroSettings        = config_inner.HeroSettings
	MapTopology         = config_inner.MapTopology
	TournamentRules     = config_inner.TournamentRules
	ZoneConfig          = config_inner.ZoneConfig
)

const (
	BonusTownPortalFree   = config_inner.BonusTownPortalFree
	BonusSpell            = config_inner.BonusSpell
	BonusUnitMultiplier   = config_inner.BonusUnitMultiplier
	BonusMovementBonus    = config_inner.BonusMovementBonus
	BonusStartingItem     = config_inner.BonusStartingItem
	BonusStartingGold     = config_inner.BonusStartingGold
	BonusStartingGems     = config_inner.BonusStartingGems
	BonusStartingCrystals = config_inner.BonusStartingCrystals
	BonusStartingMercury  = config_inner.BonusStartingMercury
	BonusStartingWood     = config_inner.BonusStartingWood
	BonusStartingOre      = config_inner.BonusStartingOre
)

const (
	TopologyRing        MapTopology = config_inner.TopologyRing
	TopologyHubAndSpoke MapTopology = config_inner.TopologyHubAndSpoke
	TopologyChain       MapTopology = config_inner.TopologyChain
	TopologySharedWeb   MapTopology = config_inner.TopologySharedWeb
	TopologyRandom      MapTopology = config_inner.TopologyRandom
	TopologyCircles     MapTopology = config_inner.TopologyCircles
	TopologySquare      MapTopology = config_inner.TopologySquare
	TopologyGeometric   MapTopology = config_inner.TopologyGeometric
	TopologyCross       MapTopology = config_inner.TopologyCross
	TopologyFractal     MapTopology = config_inner.TopologyFractal
)
