package editorStateRepository_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestWhenEditorStateRepositoryIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	repository := repositories.NewEditorStateRepository()

	// Assert
	assert.NotNil(t, repository)
}
