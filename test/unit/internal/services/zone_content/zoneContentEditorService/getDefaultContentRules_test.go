package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNoGuardedRuleIsOffered_ThereIsNoDefaultRule(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	rules := service.GetDefaultContentRules(dtos.ContentRuleEditorOptionsDto{
		Rules: []dtos.ContentRuleOptionDto{{Key: dtos.ContentRuleKeyVariant}},
	})

	// Assert
	assert.Empty(t, rules)
}

func TestWhenAGuardedRuleIsOffered_TheDefaultRuleIsNamedAfterIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()

	// Act
	rules := service.GetDefaultContentRules(dtos.ContentRuleEditorOptionsDto{
		Rules: []dtos.ContentRuleOptionDto{{Key: dtos.ContentRuleKeyGuarded, Name: name}},
	})

	// Assert
	require.Len(t, rules, 1)
	assert.Equal(t, name, rules[0].Name)
}

func TestWhenAGuardedRuleIsOffered_TheDefaultRuleIsGuarded(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	rules := service.GetDefaultContentRules(dtos.ContentRuleEditorOptionsDto{
		Rules: []dtos.ContentRuleOptionDto{{Key: dtos.ContentRuleKeyGuarded, Name: gofakeit.Word()}},
	})

	// Assert
	require.Len(t, rules, 1)
	require.NotNil(t, rules[0].IsGuarded)
	assert.True(t, *rules[0].IsGuarded)
}
