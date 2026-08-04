package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenConfigurationProvided_ReturnsNonNilGenerator(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	generator := newTemplateGenerator(configuration)

	// Assert
	assert.NotNil(t, generator)
}

func TestWhenConfigurationProvided_GeneratesFromProvidedConfiguration(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedName := gofakeit.InputName()
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = expectedName

	// Act
	generator := newTemplateGenerator(configuration)

	// Assert
	assert.Equal(t, expectedName, generator.Generate().Name)
}
