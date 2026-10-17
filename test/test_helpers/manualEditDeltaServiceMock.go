package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/mock"
)

// ManualEditDeltaServiceMock is a testify mock of
// connection_editor_interfaces.IManualEditDeltaService, used to unit-test
// collaborators without the real derivation logic.
type ManualEditDeltaServiceMock struct {
	mock.Mock
}

func (this *ManualEditDeltaServiceMock) Capture(
	request models.ManualEditCaptureRequest) editor_state_model.ManualEditSettings {
	arguments := this.Called(request)
	settings, _ := arguments.Get(0).(editor_state_model.ManualEditSettings)
	return settings
}

func (this *ManualEditDeltaServiceMock) Expand(request models.ManualEditExpandRequest) models.ManualEditExpansion {
	arguments := this.Called(request)
	expansion, _ := arguments.Get(0).(models.ManualEditExpansion)
	return expansion
}
