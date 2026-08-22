package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenUpdateChangesPlayerCount_ChangeIsApplied(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	playerCount := gofakeit.Number(3, 8)

	// Act
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.PlayerCount = playerCount })

	// Assert
	assert.Equal(t, playerCount, state.GetCurrentState().PlayerCount)
}

func TestWhenUpdateSetsPlayerCountAboveMaximum_PlayerCountIsClamped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithValidation(func(
		stateDto editor_state_model.EditorState,
		_ bool,
	) editor_state_dto.EditorStateValidationDto {
		stateDto.PlayerCount = 8
		return editor_state_dto.EditorStateValidationDto{State: stateDto}
	})
	tooManyPlayers := gofakeit.Number(9, 100)

	// Act
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.PlayerCount = tooManyPlayers })

	// Assert
	assert.Equal(t, 8, state.GetCurrentState().PlayerCount)
}

func TestWhenUpdateSetsUnknownGameMode_GameModeIsResetToClassic(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithValidation(func(
		stateDto editor_state_model.EditorState,
		_ bool,
	) editor_state_dto.EditorStateValidationDto {
		stateDto.GameMode = registry.GetGameModeValues().Classic
		return editor_state_dto.EditorStateValidationDto{State: stateDto}
	})

	// Act
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.GameMode = "NotARealGameMode" })

	// Assert
	assert.Equal(t, registry.GetGameModeValues().Classic, state.GetCurrentState().GameMode)
}
