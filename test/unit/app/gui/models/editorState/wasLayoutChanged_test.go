package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyChangedSinceSnapshot_ReportsLayoutChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.Topology = config_inner.TopologyChain })

	// Act
	layoutChanged := state.WasLayoutChanged()

	// Assert
	assert.True(t, layoutChanged)
}

func TestWhenOnlyDensityChangedSinceSnapshot_ReportsLayoutNotChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.ResourceDensityPercent = 50 })

	// Act
	layoutChanged := state.WasLayoutChanged()

	// Assert
	assert.False(t, layoutChanged)
}

// A freshly created or freshly loaded state has no previous snapshot to
// compare against, so the layout cannot have changed - and the comparison must
// not dereference the absent snapshot.
func TestWhenNoPreviousStateExists_ReportsLayoutNotChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.Topology = config_inner.TopologyChain })

	// Act
	layoutChanged := state.WasLayoutChanged()

	// Assert
	assert.False(t, layoutChanged)
}
