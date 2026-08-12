package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
)

func NewGenerationTuning(
	configuration *config.GeneratorConfig,
	totalZoneCount int,
) models.GenerationTuning {
	return generation_tuning.NewGenerationTuningFactory().Create(configuration, totalZoneCount)
}
