package tournament_variant

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IClusterService interface {
	CreateClusterVariant(
		configuration config.GeneratorConfig,
		tuning models.GenerationTuning,
		allNeutralZonePlans, playerNeutralZonePlans neutral_zone.Plans,
		playerIndex int,
		playerLabel string) ([]template_model.Zone, []template_model.Connection)
}
