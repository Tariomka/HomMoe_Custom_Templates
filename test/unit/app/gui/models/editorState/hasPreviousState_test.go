package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoSnapshotWasTaken_ReportsNoPreviousState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	hasPrevious := state.HasPreviousState()

	// Assert
	assert.False(t, hasPrevious)
}

func TestWhenSnapshotWasTaken_ReportsPreviousState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()

	// Act
	hasPrevious := state.HasPreviousState()

	// Assert
	assert.True(t, hasPrevious)
}
