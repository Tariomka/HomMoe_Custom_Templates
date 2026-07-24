package generationTuningFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactoryIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	factory := generation_tuning.NewGenerationTuningFactory()

	// Assert
	assert.NotNil(t, factory)
}
