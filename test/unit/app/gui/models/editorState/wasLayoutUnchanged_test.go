package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoSnapshotExists_DoesNotReportLayoutUnchanged(t *testing.T) {
	// Arrange
	state := models.NewEditorState()

	// Act
	layoutUnchanged := state.WasLayoutUnchanged()

	// Assert
	assert.False(t, layoutUnchanged)
}

func TestWhenPlayerCountChangedSinceSnapshot_DoesNotReportLayoutUnchanged(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	layoutUnchanged := state.WasLayoutUnchanged()

	// Assert
	assert.False(t, layoutUnchanged)
}

func TestWhenOnlyHeroCountChangedSinceSnapshot_ReportsLayoutUnchanged(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.HeroCountMax++ })

	// Act
	layoutUnchanged := state.WasLayoutUnchanged()

	// Assert
	assert.True(t, layoutUnchanged)
}
