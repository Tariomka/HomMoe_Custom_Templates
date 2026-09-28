package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenAContentRuleIsComposed_TheEmbeddedCompositionReachesTheService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	composition := fakeComposition()
	fixture.contentEditor.On("ComposeContentRule", mock.Anything).Return(editor_state_model.ContentRuleRow{}, false)

	// Act
	fixture.handler.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{ContentRuleComposition: composition})

	// Assert
	fixture.contentEditor.AssertCalled(t, "ComposeContentRule", composition)
}

func TestWhenTheServiceAcceptsTheComposition_ReturnsAValidResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	composition := fakeComposition()
	rule := editor_state_model.ContentRuleRow{Name: gofakeit.Word(), DistanceName: gofakeit.Word()}
	fixture.contentEditor.On("ComposeContentRule", composition).Return(rule, true)

	// Act
	result := fixture.handler.ComposeContentRule(
		dtos.ContentRuleCompositionRequestDto{ContentRuleComposition: composition})

	// Assert
	assert.Equal(t, dtos.ContentRuleCompositionResultDto{Rule: rule, Valid: true}, result)
}

func TestWhenTheServiceRejectsTheComposition_ReturnsAnInvalidResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	composition := fakeComposition()
	rule := editor_state_model.ContentRuleRow{Name: gofakeit.Word()}
	fixture.contentEditor.On("ComposeContentRule", composition).Return(rule, false)

	// Act
	result := fixture.handler.ComposeContentRule(
		dtos.ContentRuleCompositionRequestDto{ContentRuleComposition: composition})

	// Assert
	assert.Equal(t, dtos.ContentRuleCompositionResultDto{Rule: rule, Valid: false}, result)
}

// fakeComposition fills every field so an unwrapped pass-through is observable.
func fakeComposition() content_rule_model.ContentRuleComposition {
	return content_rule_model.ContentRuleComposition{
		Key:             content_rule_model.ContentRuleKey(gofakeit.Word()),
		Name:            gofakeit.Word(),
		DistanceNames:   []string{gofakeit.Word()},
		DistanceIndex:   gofakeit.Number(0, 10),
		IsGuarded:       gofakeit.Bool(),
		IsSoloEncounter: gofakeit.Bool(),
		VariantIDs:      []int{gofakeit.Number(1, 500)},
		VariantIndex:    gofakeit.Number(0, 10),
	}
}
