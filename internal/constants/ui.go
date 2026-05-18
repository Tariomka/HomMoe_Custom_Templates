package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

var (
	GameModes      = []string{"Classic", "SingleHero"}
	MapSizes       = []int{64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 240}
	ExpMapSizes    = []int{256, 272, 288, 304, 320, 336, 352, 368, 384, 400, 416, 432, 448, 464, 480, 496, 512}
	TopologyLabels = []string{"Random", "Ring", "Hub", "Chain", "Shared Web"}
	TopologyValues = []models.MapTopology{
		generator.TopologyRandom,
		generator.TopologyDefault,
		generator.TopologyHubAndSpoke,
		generator.TopologyChain,
		generator.TopologySharedWeb,
	}
	VictoryLabels = []string{"Standard", "Lost Starting City", "Hold City", "Tournament"}
	VictoryIDs    = []string{"win_condition_1", "win_condition_3", "win_condition_5", "win_condition_6"}
	RoadDistances = []string{"Any", "Next To", "Near", "Medium", "Far", "Very Far"}

	MainTabLabels = []string{
		"Map Setup",
		"Generation Options",
		"Game Rules",
		"Zone Content (EXP)",
	}
)
