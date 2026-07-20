package tournament_variant

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type IClusterService interface {
	CreateClusterVariant(
		configuration config.GeneratorConfig,
		tuning models.GenerationTuning,
		allNeutralZonePlans, playerNeutralZonePlans neutral_zone.Plans,
		playerIndex int,
		playerLabel string) ([]entities.Zone, []entities.Connection)
}
