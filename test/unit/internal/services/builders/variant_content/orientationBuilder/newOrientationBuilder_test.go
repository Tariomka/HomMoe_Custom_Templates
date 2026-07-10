package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyOrientation(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	builder := variant_content.NewOrientationBuilder()

	// Assert
	assert.Equal(t, entities.Orientation{}, builder.Build())
}
