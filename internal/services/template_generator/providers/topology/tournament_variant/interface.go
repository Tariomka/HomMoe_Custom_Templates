package tournament_variant

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type IClusterService interface {
	CreateClusterVariant(
		configuration config.GeneratorConfig,
		tuning models.GenerationTuning,
		allNeutralZonePlans, playerNeutralZonePlans models.NeutralZonePlans,
		playerIndex int,
		playerLabel string) ([]template.Zone, []template.Connection)
}
