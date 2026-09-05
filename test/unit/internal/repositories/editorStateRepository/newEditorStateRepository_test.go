package editorStateRepository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenEditorStateRepositoryIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	repository := newRepository()

	// Assert
	assert.NotNil(t, repository)
}
