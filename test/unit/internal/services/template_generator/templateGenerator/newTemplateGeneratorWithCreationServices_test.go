package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreationServicesAreProvided_ReturnsGenerator(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	creationServices := zones.NewCreationServices(nil, nil)

	// Act
	generator := template_generator.NewTemplateGeneratorWithCreationServices(configuration, creationServices)

	// Assert
	assert.NotNil(t, generator)
}

func TestWhenCreationServicesAreOmitted_ReturnsGenerator(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	generator := template_generator.NewTemplateGeneratorWithCreationServices(configuration, nil)

	// Assert
	assert.NotNil(t, generator)
}
