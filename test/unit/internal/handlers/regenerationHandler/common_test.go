// Package regenerationHandler_test contains shared arrangement helpers for the
// regenerationHandler.go unit tests.
package regenerationHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/regeneration"
	"github.com/stretchr/testify/mock"
)

type decisionServiceMock struct {
	mock.Mock
}

func (this *decisionServiceMock) DecideRegeneration(request regeneration.DecisionRequest) regeneration.Decision {
	args := this.Called(request)
	return args.Get(0).(regeneration.Decision)
}

func (this *decisionServiceMock) DecideManualEditReapplication(
	previous, current *editor_state_model.EditorState) regeneration.ManualEditDecision {
	args := this.Called(previous, current)
	return args.Get(0).(regeneration.ManualEditDecision)
}
