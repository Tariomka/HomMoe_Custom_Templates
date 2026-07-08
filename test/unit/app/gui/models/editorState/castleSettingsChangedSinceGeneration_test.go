package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingWasGeneratedYet_NoCastleChangesAreReported(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerZoneCastles++ })

	// Act
	changes := state.CastleSettingsChangedSinceGeneration()

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{}, changes)
}

func TestWhenPlayerCastleCountChangedSinceSnapshot_PlayerCastleChangeIsReported(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerZoneCastles++ })

	// Act
	changes := state.CastleSettingsChangedSinceGeneration()

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{PlayerCastles: true}, changes)
}

func TestWhenCastleCountsAreUnchangedSinceSnapshot_NoCastleChangesAreReported(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.ResourceDensityPercent = 50 })

	// Act
	changes := state.CastleSettingsChangedSinceGeneration()

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChanges{}, changes)
}
