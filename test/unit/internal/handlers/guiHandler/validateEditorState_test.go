package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsValid_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()

	// Act
	result := handler.ValidateEditorState(stateDto, true)

	// Assert
	assert.Empty(t, result.Warnings)
}

func TestWhenFixIssuesIsTrue_ReturnsNormalizedState(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.PlayerCount = 50

	// Act
	result := handler.ValidateEditorState(stateDto, true)

	// Assert
	assert.Equal(t, 8, result.State.PlayerCount)
}

func TestWhenFixIssuesIsFalse_ReturnsOriginalState(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.PlayerCount = 50

	// Act
	result := handler.ValidateEditorState(stateDto, false)

	// Assert
	assert.Equal(t, stateDto, result.State)
}

func TestWhenFixIssuesIsFalse_ReturnsIssueMessages(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.PlayerCount = 50

	// Act
	result := handler.ValidateEditorState(stateDto, false)

	// Assert
	assert.Equal(t, []string{"playerCount 50 is outside [2, 8]"}, result.Warnings)
}

func TestWhenStateIsInvalid_ReturnsIssueMessages(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.PlayerCount = 50

	// Act
	result := handler.ValidateEditorState(stateDto, true)

	// Assert
	assert.Equal(t, []string{"playerCount 50 is outside [2, 8]"}, result.Warnings)
}

func TestWhenStateHasMultipleIssues_ReturnsMessagesInValidationOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.PlayerCount = 50
	stateDto.NeutralZoneCount = -1

	// Act
	result := handler.ValidateEditorState(stateDto, true)

	// Assert
	assert.Equal(t, []string{
		"playerCount 50 is outside [2, 8]",
		"neutralZoneCount -1 is negative",
	}, result.Warnings)
}

func TestWhenAdvancedModeIsEnabled_ZeroesSimpleNeutralCount(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.AdvancedMode = true
	stateDto.NeutralZoneCount = gofakeit.Number(1, 10)

	// Act
	result := handler.ValidateEditorState(stateDto, true)

	// Assert
	assert.Zero(t, result.State.NeutralZoneCount)
}

func TestWhenSimpleModeIsEnabled_ZeroesAdvancedNeutralCounts(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.NeutralLowestNoCastleCount = gofakeit.Number(1, 10)
	stateDto.NeutralLowestCastleCount = gofakeit.Number(1, 10)
	stateDto.NeutralLowNoCastleCount = gofakeit.Number(1, 10)
	stateDto.NeutralLowCastleCount = gofakeit.Number(1, 10)
	stateDto.NeutralMediumNoCastleCount = gofakeit.Number(1, 10)
	stateDto.NeutralMediumCastleCount = gofakeit.Number(1, 10)
	stateDto.NeutralHighNoCastleCount = gofakeit.Number(1, 10)
	stateDto.NeutralHighCastleCount = gofakeit.Number(1, 10)
	expected := dtos.NewDefaultEditorStateDto()

	// Act
	result := handler.ValidateEditorState(stateDto, true)

	// Assert
	assert.Equal(t, expected, result.State)
}

func TestWhenFixIssuesIsFalse_PreservesInactiveNeutralCounts(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()
	stateDto.NeutralLowNoCastleCount = gofakeit.Number(1, 10)

	// Act
	result := handler.ValidateEditorState(stateDto, false)

	// Assert
	assert.Equal(t, stateDto, result.State)
}
