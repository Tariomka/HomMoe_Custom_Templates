package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsValid_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	stateDto := dtos.NewDefaultEditorStateDto()

	// Act
	result := handler.ValidateEditorState(stateDto, true)

	// Assert
	assert.Empty(t, result.Warnings)
}

func TestWhenFixIssuesIsTrue_ReturnsNormalizedState(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
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
	handler := handlers.NewGuiHandler()
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
	handler := handlers.NewGuiHandler()
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
	handler := handlers.NewGuiHandler()
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
	handler := handlers.NewGuiHandler()
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
