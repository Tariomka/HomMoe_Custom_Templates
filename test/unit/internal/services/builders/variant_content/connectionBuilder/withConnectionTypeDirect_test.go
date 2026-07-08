package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenDirectTypeIsChosen_SetsDirectConnectionTypeOnBuiltConnection(t *testing.T) {
	// Arrange
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithConnectionTypeDirect().Build()

	// Assert
	assert.Equal(t, entities.Connection{ConnectionType: "Direct"}, connection)
}
