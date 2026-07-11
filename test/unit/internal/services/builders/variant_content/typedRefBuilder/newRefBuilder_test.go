package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyReference(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	builder := variant_content.NewRefBuilder()

	// Assert
	assert.Equal(t, entities.TypedRef{}, builder.Build())
}
