package v2EditorStateRepository_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v2"
	"github.com/stretchr/testify/assert"
)

// Writing a v2 file would produce something this build could not round-trip,
// so the repository refuses rather than quietly emitting the old shape.
func TestWhenAV2StateIsSaved_TheWriteIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := newRepository()

	// Act
	_, err := repository.Save(t.TempDir(), "State", editor_state_v2.EditorState{})

	// Assert
	assert.Error(t, err)
}

func TestWhenAV2StateIsSaved_NoPathIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := newRepository()

	// Act
	writtenPath, _ := repository.Save(t.TempDir(), "State", editor_state_v2.EditorState{})

	// Assert
	assert.Empty(t, writtenPath)
}
