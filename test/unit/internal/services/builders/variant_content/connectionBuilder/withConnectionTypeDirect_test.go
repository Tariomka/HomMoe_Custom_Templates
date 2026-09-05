package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenDirectTypeIsChosen_SetsDirectConnectionTypeOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithConnectionTypeDirect().Build()

	// Assert
	assert.Equal(t, template_model.Connection{ConnectionType: "Direct"}, connection)
}
