package legacyEditorStateRepository_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
	"github.com/stretchr/testify/assert"
)

// Writing a v1 file would produce something this build could not round-trip, so
// the legacy repository refuses rather than quietly emitting the old shape.
func TestWhenALegacyStateIsSaved_TheWriteIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := newRepository()

	// Act
	_, err := repository.Save(t.TempDir(), "State", editor_state_v1.EditorState{})

	// Assert
	assert.Error(t, err)
}

func TestWhenALegacyStateIsSaved_NoPathIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := newRepository()

	// Act
	writtenPath, _ := repository.Save(t.TempDir(), "State", editor_state_v1.EditorState{})

	// Assert
	assert.Empty(t, writtenPath)
}
