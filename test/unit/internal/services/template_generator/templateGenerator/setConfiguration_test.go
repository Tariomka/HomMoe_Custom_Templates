package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNewConfigurationProvided_GeneratesWithNewConfiguration(t *testing.T) {
	// Arrange
	expectedName := gofakeit.InputName()
	generator := template_generator.NewTemplateGenerator(config.NewGeneratorConfig())
	newConfiguration := config.NewGeneratorConfig()
	newConfiguration.TemplateName = expectedName

	// Act
	generator.SetConfiguration(newConfiguration)

	// Assert
	assert.Equal(t, expectedName, generator.Generate().Name)
}

func TestWhenNilConfigurationProvided_KeepsPreviousConfiguration(t *testing.T) {
	// Arrange
	expectedName := gofakeit.InputName()
	originalConfiguration := config.NewGeneratorConfig()
	originalConfiguration.TemplateName = expectedName
	generator := template_generator.NewTemplateGenerator(originalConfiguration)

	// Act
	generator.SetConfiguration(nil)

	// Assert
	assert.Equal(t, expectedName, generator.Generate().Name)
}
