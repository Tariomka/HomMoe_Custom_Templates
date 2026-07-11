package mandatoryContentBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRulesCallbackIsProvided_AppendsCallbackRulesOnBuiltItem(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	expectedRules := []entities.PlacementRule{
		{Type: "Road", Weight: gofakeit.Number(1, 100)},
		{Type: "MainObject", Weight: gofakeit.Number(1, 100)},
	}
	builder := mandatory_content.NewContentBuilder(expectedSid)

	// Act
	item := builder.WithRulesCallback(func() []entities.PlacementRule { return expectedRules }).Build()

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{SID: expectedSid, Rules: expectedRules}, item)
}
