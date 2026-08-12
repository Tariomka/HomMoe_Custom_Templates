package previewRepository_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestWhenPreviewRepositoryIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	repository := repositories.NewPreviewRepository()

	// Assert
	assert.NotNil(t, repository)
}
