package templateRepository_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestWhenTemplateRepositoryIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	repository := repositories.NewTemplateRepository()

	// Assert
	assert.NotNil(t, repository)
}
