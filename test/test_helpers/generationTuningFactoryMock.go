package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/mock"
)

// GenerationTuningFactoryMock is a testify mock of
// generation_tuning.IGenerationTuningFactory, used to unit-test collaborators
// with a fixed tuning instead of the derived one.
type GenerationTuningFactoryMock struct {
	mock.Mock
}

func (this *GenerationTuningFactoryMock) Create(
	configuration *config.GeneratorConfig,
	totalZoneCount int) models.GenerationTuning {
	arguments := this.Called(configuration, totalZoneCount)
	tuning, _ := arguments.Get(0).(models.GenerationTuning)
	return tuning
}
