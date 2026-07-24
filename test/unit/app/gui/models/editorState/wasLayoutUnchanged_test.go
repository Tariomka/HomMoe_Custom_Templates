package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNoSnapshotExists_DoesNotReportLayoutUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	layoutUnchanged := state.WasLayoutUnchanged()

	// Assert
	assert.False(t, layoutUnchanged)
}

func TestWhenPlayerCountChangedSinceSnapshot_DoesNotReportLayoutUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	layoutUnchanged := state.WasLayoutUnchanged()

	// Assert
	assert.False(t, layoutUnchanged)
}

func TestWhenOnlyHeroCountChangedSinceSnapshot_ReportsLayoutUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.HeroCountMax++ })

	// Act
	layoutUnchanged := state.WasLayoutUnchanged()

	// Assert
	assert.True(t, layoutUnchanged)
}
