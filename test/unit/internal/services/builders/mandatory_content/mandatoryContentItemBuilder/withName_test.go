package mandatoryContentItemBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsProvided_SetsNameOnBuiltItem(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	expectedName := gofakeit.Word()
	builder := mandatory_content.NewContentItemBuilder(expectedSid)

	// Act
	item := builder.WithName(expectedName).Build()

	// Assert
	assert.Equal(t, template_model.MandatoryContentItem{SID: expectedSid, Name: expectedName}, item)
}
