package generation_tuning

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type IGenerationTuningFactory interface {
	Create(configuration *config.GeneratorConfig, totalZoneCount int) models.GenerationTuning
}
