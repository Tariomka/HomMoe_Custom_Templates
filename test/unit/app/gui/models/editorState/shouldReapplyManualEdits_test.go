package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/stretchr/testify/assert"
)

func stateWithManualEdits() *models.EditorState {
	state := models.NewEditorState()
	state.SetManualEdits([]entities.Zone{{Name: "Zone A"}}, nil)
	return state
}

func TestWhenNoManualEditsExist_ReapplyIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	shouldReapply := state.ShouldReapplyManualEdits()

	// Assert
	assert.False(t, shouldReapply)
}

func TestWhenManualEditsExistWithoutSnapshot_ReapplyIsAllowed(t *testing.T) {
	t.Parallel()
	// Arrange
	state := stateWithManualEdits()

	// Act
	shouldReapply := state.ShouldReapplyManualEdits()

	// Assert
	assert.True(t, shouldReapply)
}

func TestWhenLayoutChangedAfterManualEdits_ReapplyIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	state := stateWithManualEdits()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.Topology = config_inner.TopologyChain })

	// Act
	shouldReapply := state.ShouldReapplyManualEdits()

	// Assert
	assert.False(t, shouldReapply)
}

func TestWhenOnlyNonLayoutOptionsChangedAfterManualEdits_ReapplyIsAllowed(t *testing.T) {
	t.Parallel()
	// Arrange
	state := stateWithManualEdits()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.ResourceDensityPercent = 50 })

	// Act
	shouldReapply := state.ShouldReapplyManualEdits()

	// Assert
	assert.True(t, shouldReapply)
}
