package regenerationDecisionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/regeneration"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoPreviousGenerationAndManualEditsExist_Reapplies(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(nil, stateWithManualEdits())

	// Assert
	assert.Equal(t, regeneration.ManualEditDecision{
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
	assert.Equal(t, regeneration.ManualEditDecision{}, decision)
}

func TestWhenNoManualEditsExist_DoesNotReapply(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), defaultState())

	// Assert
	assert.Equal(t, regeneration.ManualEditDecision{}, decision)
}

func TestWhenManualEditsExistAndLayoutUnchanged_Reapplies(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), stateWithManualEdits())

	// Assert
	assert.Equal(t, regeneration.ManualEditDecision{
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
	current.ManualZones = manualZones()

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

// Manual edits that survived into the generation of an effective-mode change
// belong to the map that mode replaced, so they are not reapplied.
func TestWhenManualEditsExistButTournamentModeWasEntered_DoesNotReapply(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := tournamentState()
	current.ManualZones = manualZones()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.Nil(t, decision.ReapplyWithCastleChanges)
}

func TestWhenManualEditsExistButArenaModeWasEntered_DoesNotReapply(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := arenaState()
	current.ManualZones = manualZones()

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.Nil(t, decision.ReapplyWithCastleChanges)
}

func TestWhenManualEditsExistWithAStaleGeneratorVersion_ReportsGeneratorVersionChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := stateWithManualEdits()
	current.GeneratorVersion = common_generator.GeneratorOutputVersion + 1

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.True(t, decision.GeneratorVersionChanged)
}

func TestWhenManualEditsExistWithAStaleGeneratorVersionAndNoPreviousGeneration_ReportsGeneratorVersionChanged(
	t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := stateWithManualEdits()
	current.GeneratorVersion = common_generator.GeneratorOutputVersion + 1

	// Act
	decision := service.DecideManualEditReapplication(nil, current)

	// Assert
	assert.True(t, decision.GeneratorVersionChanged)
}

func TestWhenManualEditsExistWithTheCurrentGeneratorVersion_ReportsNoGeneratorVersionChange(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := stateWithManualEdits()
	current.GeneratorVersion = common_generator.GeneratorOutputVersion

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.False(t, decision.GeneratorVersionChanged)
}

func TestWhenManualEditsExistWithTheCurrentGeneratorVersionAndNoPreviousGeneration_ReportsNoGeneratorVersionChange(
	t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := stateWithManualEdits()
	current.GeneratorVersion = common_generator.GeneratorOutputVersion

	// Act
	decision := service.DecideManualEditReapplication(nil, current)

	// Assert
	assert.False(t, decision.GeneratorVersionChanged)
}

func TestWhenManualEditsExistWithoutAGeneratorVersion_ReportsNoGeneratorVersionChange(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := stateWithManualEdits()
	current.GeneratorVersion = 0

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.False(t, decision.GeneratorVersionChanged)
}

func TestWhenThereAreNoManualEditsButAStaleGeneratorVersion_ReportsNoGeneratorVersionChange(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := defaultState()
	current.GeneratorVersion = common_generator.GeneratorOutputVersion + 1

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.False(t, decision.GeneratorVersionChanged)
}

func TestWhenTheLayoutChangedWithAStaleGeneratorVersion_ReportsNoGeneratorVersionChange(t *testing.T) {
	t.Parallel()
	// Arrange
	service := editor.NewRegenerationDecisionService()
	current := layoutChangedState()
	current.ManualZones = manualZones()
	current.GeneratorVersion = common_generator.GeneratorOutputVersion + 1

	// Act
	decision := service.DecideManualEditReapplication(defaultState(), current)

	// Assert
	assert.False(t, decision.GeneratorVersionChanged)
}

func manualZones() []editor_state_model.ManualZoneEdit {
	return editor_state_model.ManualEditSettingsFromSnapshot(
		[]template_model.Zone{{Name: gofakeit.Word()}}, nil).ManualZones
}

func stateWithManualEdits() *editor_state_model.EditorState {
	state := defaultState()
	state.ManualZones = manualZones()
	return state
}
