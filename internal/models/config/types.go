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
	TopologyDefault     MapTopology = config_inner.TopologyDefault
	TopologyHubAndSpoke MapTopology = config_inner.TopologyHubAndSpoke
	TopologyChain       MapTopology = config_inner.TopologyChain
	TopologySharedWeb   MapTopology = config_inner.TopologySharedWeb
	TopologyRandom      MapTopology = config_inner.TopologyRandom
	TopologyBalanced    MapTopology = config_inner.TopologyBalanced
)
