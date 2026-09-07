package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyVariant(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	builder := variant_content.NewVariantBuilder()

	// Assert
	assert.Equal(t, template_model.Variant{}, builder.Build())
}
