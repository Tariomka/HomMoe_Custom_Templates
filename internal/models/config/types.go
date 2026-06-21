package config

import "github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"

type (
	AdvancedSettings    = config_inner.AdvancedSettings
	BonusEntry          = config_inner.BonusEntry
	GameEndConditions   = config_inner.GameEndConditions
	GladiatorArenaRules = config_inner.GladiatorArenaRules
	HeroSettings        = config_inner.HeroSettings
	MapTopology         = config_inner.MapTopology
	TournamentRules     = config_inner.TournamentRules
	ZoneConfig          = config_inner.ZoneConfig
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

var (
	ParseBonusesJSON = config_inner.DeserializeBonusEntries
)
