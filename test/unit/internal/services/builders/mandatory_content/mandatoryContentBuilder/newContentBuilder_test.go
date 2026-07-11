package mandatoryContentBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesItemWithOnlyProvidedSid(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()

	// Act
	builder := mandatory_content.NewContentBuilder(expectedSid)

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{SID: expectedSid}, builder.Build())
}
