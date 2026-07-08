package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyBorder(t *testing.T) {
	// Arrange & Act
	builder := variant_content.NewBorderBuilder()

	// Assert
	assert.Equal(t, entities.Border{}, builder.Build())
}
