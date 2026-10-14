package regenerationHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/regeneration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenTheServiceReportsAChangedGeneratorVersion_DtoForwardsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &decisionServiceMock{}
	service.On("DecideManualEditReapplication", mock.Anything, mock.Anything).Return(regeneration.ManualEditDecision{
		ReapplyWithCastleChanges: &editor_state_model.CastleSettingChanges{},
		GeneratorVersionChanged:  true,
	})
	handler := handlers.NewRegenerationHandler(service)

	// Act
	decision := handler.DecideManualEditReapplication(
		&editor_state_dto.EditorStateDto{},
		&editor_state_dto.EditorStateDto{},
	)

	// Assert
	assert.True(t, decision.GeneratorVersionChanged)
}

func TestWhenTheServiceReportsAnUnchangedGeneratorVersion_DtoForwardsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &decisionServiceMock{}
	service.On("DecideManualEditReapplication", mock.Anything, mock.Anything).Return(regeneration.ManualEditDecision{
		ReapplyWithCastleChanges: &editor_state_model.CastleSettingChanges{},
	})
	handler := handlers.NewRegenerationHandler(service)

	// Act
	decision := handler.DecideManualEditReapplication(
		&editor_state_dto.EditorStateDto{},
		&editor_state_dto.EditorStateDto{},
	)

	// Assert
	assert.False(t, decision.GeneratorVersionChanged)
}

func TestWhenTheServiceDecidesToReapply_DtoForwardsTheCastleChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &decisionServiceMock{}
	changes := &editor_state_model.CastleSettingChanges{PlayerCastles: true}
	service.On("DecideManualEditReapplication", mock.Anything, mock.Anything).Return(regeneration.ManualEditDecision{
		ReapplyWithCastleChanges: changes,
	})
	handler := handlers.NewRegenerationHandler(service)

	// Act
	decision := handler.DecideManualEditReapplication(
		&editor_state_dto.EditorStateDto{},
		&editor_state_dto.EditorStateDto{},
	)

	// Assert
	assert.Equal(t, dtos.ManualEditDecisionDto{ReapplyWithCastleChanges: changes}, decision)
}

func TestWhenTheServiceDecidesToDropTheEdits_DtoIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &decisionServiceMock{}
	service.On("DecideManualEditReapplication", mock.Anything, mock.Anything).Return(regeneration.ManualEditDecision{
		GeneratorVersionChanged: true,
	})
	handler := handlers.NewRegenerationHandler(service)

	// Act
	decision := handler.DecideManualEditReapplication(
		&editor_state_dto.EditorStateDto{},
		&editor_state_dto.EditorStateDto{},
	)

	// Assert
	assert.Equal(t, dtos.ManualEditDecisionDto{}, decision)
}

func TestWhenThePreviousStateIsNil_ServiceReceivesNilPrevious(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &decisionServiceMock{}
	current := &editor_state_dto.EditorStateDto{}
	service.On("DecideManualEditReapplication",
		(*editor_state_model.EditorState)(nil), &current.EditorState).
		Return(regeneration.ManualEditDecision{}).Once()
	handler := handlers.NewRegenerationHandler(service)

	// Act
	handler.DecideManualEditReapplication(nil, current)

	// Assert
	service.AssertExpectations(t)
}

func TestWhenBothStatesAreGiven_ServiceReceivesTheirModels(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &decisionServiceMock{}
	previous := &editor_state_dto.EditorStateDto{}
	current := &editor_state_dto.EditorStateDto{}
	service.On("DecideManualEditReapplication", &previous.EditorState, &current.EditorState).
		Return(regeneration.ManualEditDecision{}).Once()
	handler := handlers.NewRegenerationHandler(service)

	// Act
	handler.DecideManualEditReapplication(previous, current)

	// Assert
	service.AssertExpectations(t)
}
