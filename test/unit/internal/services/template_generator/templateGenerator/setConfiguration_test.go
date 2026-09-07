package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNewConfigurationProvided_GeneratesWithNewConfiguration(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedName := gofakeit.InputName()
	generator := test_helpers.NewTemplateGenerator(config.NewGeneratorConfig())
	newConfiguration := config.NewGeneratorConfig()
	newConfiguration.TemplateName = expectedName

	// Act
	generator.SetConfiguration(newConfiguration)

	// Assert
	template, _ := generateTemplate(generator)
	assert.Equal(t, expectedName, template.Name)
}

func TestWhenNilConfigurationProvided_KeepsPreviousConfiguration(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedName := gofakeit.InputName()
	originalConfiguration := config.NewGeneratorConfig()
	originalConfiguration.TemplateName = expectedName
	generator := test_helpers.NewTemplateGenerator(originalConfiguration)

	// Act
	generator.SetConfiguration(nil)

	// Assert
	template, _ := generateTemplate(generator)
	assert.Equal(t, expectedName, template.Name)
}
