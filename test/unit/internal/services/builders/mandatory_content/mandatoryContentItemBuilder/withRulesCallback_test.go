package mandatoryContentItemBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRulesCallbackIsProvided_AppendsCallbackRulesOnBuiltItem(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	expectedRules := []template_model.PlacementRule{
		{Type: "Road", Weight: gofakeit.Number(1, 100)},
		{Type: "MainObject", Weight: gofakeit.Number(1, 100)},
	}
	builder := mandatory_content.NewContentItemBuilder(expectedSid)

	// Act
	item := builder.WithRulesCallback(func() []template_model.PlacementRule { return expectedRules }).Build()

	// Assert
	assert.Equal(t, template_model.MandatoryContentItem{SID: expectedSid, Rules: expectedRules}, item)
}
