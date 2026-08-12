package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/mock"
)

// ManualReapplyServiceMock is a testify mock of
// connection_editor.IManualReapplyService, used to unit-test collaborators
// without the real castle-propagation logic.
type ManualReapplyServiceMock struct {
	mock.Mock
}

func (this *ManualReapplyServiceMock) ApplyCastleSettingChanges(
	zones []entities.Zone,
	changes editor_state_dto.CastleSettingChanges,
	configuration *config.GeneratorConfig) {
	this.Called(zones, changes, configuration)
}

func (this *ManualReapplyServiceMock) SetNeutralZoneCastleCount(
	zone *entities.Zone,
	castleCount int,
	tuning models.GenerationTuning) {
	this.Called(zone, castleCount, tuning)
}
