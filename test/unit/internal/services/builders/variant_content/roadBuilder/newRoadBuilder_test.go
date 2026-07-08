package roadBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyRoad(t *testing.T) {
	// Arrange & Act
	builder := variant_content.NewRoadBuilder()

	// Assert
	assert.Equal(t, entities.Road{}, builder.Build())
}
