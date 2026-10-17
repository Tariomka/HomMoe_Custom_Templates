package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// ManualReapplyServiceMock is a testify mock of
// connection_editor_interfaces.IManualReapplyService, used to unit-test collaborators
// without the real castle-propagation logic.
type ManualReapplyServiceMock struct {
	mock.Mock
}

func (this *ManualReapplyServiceMock) ApplyCastleSettingChanges(
	zones []template_model.Zone,
	eligible []bool,
	changes editor_state_model.CastleSettingChanges,
	configuration *config.GeneratorConfig) {
	this.Called(zones, eligible, changes, configuration)
}

func (this *ManualReapplyServiceMock) SetNeutralZoneCastleCount(
	zone *template_model.Zone,
	castleCount int,
	tuning models.GenerationTuning) {
	this.Called(zone, castleCount, tuning)
}
