package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenSnapshotIsTaken_PreviousStateExists(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	state.SnapshotCurrentState()

	// Assert
	assert.True(t, state.HasPreviousState())
}

func TestWhenCurrentStateChangesAfterSnapshot_SnapshotKeepsOldValues(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()

	// Act
	state.UpdateCurrentState(func(dto *editor_state_dto.EditorStateDto) { dto.PlayerCount++ })

	// Assert - the snapshot still holds the old player count, so the state reads as changed
	assert.True(t, state.WasStateChanged())
}

// TestWhenSnapshotTakenAndContentRowMutatedInPlace_ReportsStateChanged locks in
// the deep snapshot: a shallow struct copy would share the content-row backing
// array, so an in-place edit would leave the state looking unchanged and the
// editor would neither mark the file dirty nor regenerate.
func TestWhenSnapshotTakenAndContentRowMutatedInPlace_ReportsStateChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *editor_state_dto.EditorStateDto) {
		dto.PlayerZoneContentRows = []models.ZoneContentRowSave{{Sid: "sawmill", Count: 1}}
	})
	state.SnapshotCurrentState()

	// Act
	state.UpdateCurrentState(func(dto *editor_state_dto.EditorStateDto) { dto.PlayerZoneContentRows[0].Count = 5 })

	// Assert
	assert.True(t, state.WasStateChanged())
}
