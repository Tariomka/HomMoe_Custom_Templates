package tournament_variant

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

type BalancedClusterService struct {
	base.TopologyBase
}

func NewBalancedClusterService() *BalancedClusterService {
	return &BalancedClusterService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *BalancedClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans models.NeutralZonePlans,
	playerIndex int,
	playerLabel string) ([]template.Zone, []template.Connection) {
	return nil, nil
}
