package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
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
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount = playerCount })

	// Assert
	assert.Equal(t, playerCount, state.GetCurrentState().PlayerCount)
}

func TestWhenUpdateEnablesAdvancedMode_SimpleNeutralZoneCountIsZeroed(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) {
		dto.AdvancedMode = true
		dto.NeutralZoneCount = gofakeit.Number(1, 10)
	})

	// Assert
	assert.Equal(t, 0, state.GetCurrentState().NeutralZoneCount)
}

func TestWhenUpdateStaysInSimpleMode_AdvancedNeutralCountsAreZeroed(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) {
		dto.NeutralLowNoCastleCount = gofakeit.Number(1, 5)
		dto.NeutralLowCastleCount = gofakeit.Number(1, 5)
		dto.NeutralMediumNoCastleCount = gofakeit.Number(1, 5)
		dto.NeutralMediumCastleCount = gofakeit.Number(1, 5)
		dto.NeutralHighNoCastleCount = gofakeit.Number(1, 5)
		dto.NeutralHighCastleCount = gofakeit.Number(1, 5)
	})

	// Assert
	assert.Equal(t, dtos.NewDefaultEditorStateDto(), state.GetCurrentState())
}

func TestWhenUpdateSetsPlayerCountAboveMaximum_PlayerCountIsClamped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithValidation(func(
		stateDto dtos.EditorStateDto,
		_ bool,
	) dtos.EditorStateValidationDto {
		stateDto.PlayerCount = 8
		return dtos.EditorStateValidationDto{State: stateDto}
	})
	tooManyPlayers := gofakeit.Number(9, 100)

	// Act
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount = tooManyPlayers })

	// Assert
	assert.Equal(t, 8, state.GetCurrentState().PlayerCount)
}

func TestWhenUpdateSetsUnknownGameMode_GameModeIsResetToClassic(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithValidation(func(
		stateDto dtos.EditorStateDto,
		_ bool,
	) dtos.EditorStateValidationDto {
		stateDto.GameMode = registry.GetGameModeValues().Classic
		return dtos.EditorStateValidationDto{State: stateDto}
	})

	// Act
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.GameMode = "NotARealGameMode" })

	// Assert
	assert.Equal(t, registry.GetGameModeValues().Classic, state.GetCurrentState().GameMode)
}
