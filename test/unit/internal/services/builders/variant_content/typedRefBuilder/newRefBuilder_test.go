package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyReference(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	builder := variant_content.NewRefBuilder()

	// Assert
	assert.Equal(t, template_model.TypedRef{}, builder.Build())
}
