package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyConnection(t *testing.T) {
	// Arrange & Act
	builder := variant_content.NewConnectionBuilder()

	// Assert
	assert.Equal(t, entities.Connection{}, builder.Build())
}
