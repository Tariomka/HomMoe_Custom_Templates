package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMandatoryContentNamesAreProvided_SetsMandatoryContentOnBuiltZone(t *testing.T) {
	// Arrange
	firstName := gofakeit.Word()
	secondName := gofakeit.Word()
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithMandatoryContent(firstName, secondName).Build()

	// Assert
	assert.Equal(t, entities.Zone{MandatoryContent: entities.StringList{firstName, secondName}}, zone)
}
