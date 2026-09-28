package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenDefaultRulesAreRequested_OnlyKeyAndNameReachTheService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	guarded := content_rule_model.ContentRuleOption{
		Key:  content_rule_model.ContentRuleKeyGuarded,
		Name: gofakeit.Word(),
	}
	variant := content_rule_model.ContentRuleOption{
		Key:  content_rule_model.ContentRuleKeyVariant,
		Name: gofakeit.Word(),
	}
	fixture.contentRules.GetContentRuleEditorOptionsFunc = func(models.SidMapping) dtos.ContentRuleEditorOptionsDto {
		return dtos.ContentRuleEditorOptionsDto{
			Rules:     []dtos.ContentRuleOptionDto{fakeOptionDto(guarded), fakeOptionDto(variant)},
			Distances: []string{gofakeit.Word()},
		}
	}
	fixture.contentEditor.On("GetDefaultContentRules", mock.Anything).Return([]editor_state_model.ContentRuleRow{})

	// Act
	fixture.handler.GetDefaultContentRules(models.SidMapping{Name: gofakeit.Word()})

	// Assert
	fixture.contentEditor.AssertCalled(
		t,
		"GetDefaultContentRules",
		[]content_rule_model.ContentRuleOption{guarded, variant},
	)
}

func TestWhenDefaultRulesAreRequested_ReturnsTheServiceResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	option := content_rule_model.ContentRuleOption{Key: content_rule_model.ContentRuleKeyGuarded, Name: gofakeit.Word()}
	expected := []editor_state_model.ContentRuleRow{{Name: gofakeit.Word()}}
	fixture.contentRules.GetContentRuleEditorOptionsFunc = func(models.SidMapping) dtos.ContentRuleEditorOptionsDto {
		return dtos.ContentRuleEditorOptionsDto{Rules: []dtos.ContentRuleOptionDto{fakeOptionDto(option)}}
	}
	fixture.contentEditor.On("GetDefaultContentRules", []content_rule_model.ContentRuleOption{option}).Return(expected)

	// Act
	rules := fixture.handler.GetDefaultContentRules(models.SidMapping{Name: gofakeit.Word()})

	// Assert
	assert.Equal(t, expected, rules)
}

// fakeOptionDto wraps option in a DTO whose presentation fields are all filled.
func fakeOptionDto(option content_rule_model.ContentRuleOption) dtos.ContentRuleOptionDto {
	return dtos.ContentRuleOptionDto{
		ContentRuleOption: option,
		Description:       gofakeit.Sentence(3),
		Marker:            gofakeit.Letter(),
		EditorKind:        content_rule_model.ContentRuleEditorKind(gofakeit.Word()),
		EditorLabel:       gofakeit.Word(),
	}
}
