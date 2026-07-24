package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNothingWasGeneratedYet_NoCastleChangesAreReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerZoneCastles++ })

	// Act
	changes := state.CastleSettingsChangedSinceGeneration()

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{}, changes)
}

func TestWhenPlayerCastleCountChangedSinceSnapshot_PlayerCastleChangeIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerZoneCastles++ })

	// Act
	changes := state.CastleSettingsChangedSinceGeneration()

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{PlayerCastles: true}, changes)
}

func TestWhenCastleCountsAreUnchangedSinceSnapshot_NoCastleChangesAreReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.ResourceDensityPercent = 50 })

	// Act
	changes := state.CastleSettingsChangedSinceGeneration()

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{}, changes)
}
