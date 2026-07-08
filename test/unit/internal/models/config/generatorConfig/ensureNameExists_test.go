package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsEmpty_SetsDefaultName(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = ""

	// Act
	configuration.EnsureNameExists()

	// Assert
	assert.Equal(t, "Custom Template", configuration.TemplateName)
}

func TestWhenNameIsAlreadySet_KeepsExistingName(t *testing.T) {
	// Arrange
	existingName := gofakeit.ProductName()
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = existingName

	// Act
	configuration.EnsureNameExists()

	// Assert
	assert.Equal(t, existingName, configuration.TemplateName)
}
