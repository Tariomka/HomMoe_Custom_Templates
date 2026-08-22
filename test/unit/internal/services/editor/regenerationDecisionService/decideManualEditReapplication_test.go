package regenerationDecisionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func manualZoneSaves() []editor_state.ManualZoneSave {
	return []editor_state.ManualZoneSave{{Zone: entities.Zone{Name: gofakeit.Word()}}}
}

func stateWithManualEdits() *editor_state_dto.EditorStateDto {
	state := defaultState()
	state.ManualZones = manualZoneSaves()
	return state
}

func TestWhenNoPreviousGenerationAndManualEditsExist_Reapplies(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(nil, stateWithManualEdits())

	// Assert
	assert.Equal(t, dtos.ManualEditDecisionDto{
		ReapplyWithCastleChanges: &editor_state_model.CastleSettingChanges{},
	}, decision)
}

func TestWhenNoPreviousGenerationAndNoManualEdits_DoesNotReapply(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(nil, defaultState())

	// Assert
	assert.Equal(t, dtos.ManualEditDecisionDto{}, decision)
}

func TestWhenNoManualEditsExist_DoesNotReapply(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), defaultState())

	// Assert
	assert.Equal(t, dtos.ManualEditDecisionDto{}, decision)
}

func TestWhenManualEditsExistAndLayoutUnchanged_Reapplies(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), stateWithManualEdits())

	// Assert
	assert.Equal(t, dtos.ManualEditDecisionDto{
		ReapplyWithCastleChanges: &editor_state_model.CastleSettingChanges{},
	}, decision)
}

// A layout-defining change invalidates the hand-made layout, so the edits are
// dropped rather than reapplied to a differently shaped map.
func TestWhenManualEditsExistButLayoutChanged_DoesNotReapply(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := layoutChangedState()
	current.ManualZones = manualZoneSaves()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.Nil(t, decision.ReapplyWithCastleChanges)
}

func TestWhenCastleOptionChangedSinceGeneration_ReportsCastleChange(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	previous := defaultState()
	current := stateWithManualEdits()
	current.PlayerZoneCastles = previous.PlayerZoneCastles + 1

	// Act
	decision := service.DecideManualEditReapplication(previous, current)

	// Assert
	assert.Equal(t,
		&editor_state_model.CastleSettingChanges{PlayerCastles: true},
		decision.ReapplyWithCastleChanges)
}

func TestWhenCastleOptionsUnchangedSinceGeneration_ReportsNoCastleChange(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), stateWithManualEdits())

	// Assert
	assert.Equal(t, &editor_state_model.CastleSettingChanges{}, decision.ReapplyWithCastleChanges)
}
