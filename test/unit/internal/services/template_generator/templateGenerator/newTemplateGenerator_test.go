package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenConfigurationProvided_ReturnsNonNilGenerator(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	generator := template_generator.NewTemplateGenerator(configuration)

	// Assert
	assert.NotNil(t, generator)
}

func TestWhenConfigurationIsNil_FallsBackToDefaultConfiguration(t *testing.T) {
	// Arrange & Act
	generator := template_generator.NewTemplateGenerator(nil)

	// Assert
	generated := generator.Generate()
	assert.Equal(t, "Custom Template", generated.Name)
}

func TestWhenConfigurationProvided_GeneratesFromProvidedConfiguration(t *testing.T) {
	// Arrange
	expectedName := gofakeit.InputName()
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = expectedName

	// Act
	generator := template_generator.NewTemplateGenerator(configuration)

	// Assert
	assert.Equal(t, expectedName, generator.Generate().Name)
}
