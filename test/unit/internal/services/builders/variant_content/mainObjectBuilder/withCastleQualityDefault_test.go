package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultCastleQualityIsChosen_SetsDefaultConstructionSidOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithCastleQualityDefault().Build()

	// Assert
	assert.Equal(t, template_model.MainObject{BuildingsConstructionSid: "default_buildings_construction"}, mainObject)
}
