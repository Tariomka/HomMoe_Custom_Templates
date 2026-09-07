package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/editor_state_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRowIsCloned_ScalarFieldsAreCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state.ZoneContentRow{Sid: gofakeit.LetterN(10), Count: gofakeit.IntRange(1, 9), IsGroup: true}

	// Act
	clone := editor_state_helpers.CloneZoneContentRow(row)

	// Assert
	assert.Equal(t, row, clone)
}

func TestWhenRowIsCloned_CopiedFieldsAreNotTheSameReferences(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state.ZoneContentRow{
		Sid:     gofakeit.LetterN(10),
		Count:   gofakeit.IntRange(1, 9),
		IsGroup: true,
		Rules: []editor_state.ContentRuleRow{
			{Name: gofakeit.LetterN(8)},
			{Name: gofakeit.LetterN(8)},
			{Name: gofakeit.LetterN(8)},
		},
	}

	// Act
	clone := editor_state_helpers.CloneZoneContentRow(row)

	// Assert
	assert.NotSame(t, &row, &clone)
	for ruleIndex := range row.Rules {
		assert.NotSame(t, &row.Rules[ruleIndex], &clone.Rules[ruleIndex])
	}
}

func TestWhenRulesAreNil_CloneRulesStayNil(t *testing.T) {
	t.Parallel()
	// Arrange - a nil slice must not become an empty one, because the editor's
	// change detection distinguishes the two.
	row := editor_state.ZoneContentRow{Sid: gofakeit.LetterN(10), Count: 1}

	// Act
	clone := editor_state_helpers.CloneZoneContentRow(row)

	// Assert
	assert.Nil(t, clone.Rules)
}

func TestWhenRuleIsReplacedInPlaceOnTheClone_SourceRuleIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	ruleName := gofakeit.LetterN(8)
	row := editor_state.ZoneContentRow{
		Sid:   gofakeit.LetterN(10),
		Count: 1,
		Rules: []editor_state.ContentRuleRow{{Name: ruleName}},
	}
	clone := editor_state_helpers.CloneZoneContentRow(row)

	// Act
	clone.Rules[0].Name = ruleName + "-changed"

	// Assert
	assert.Equal(t, ruleName, row.Rules[0].Name)
}

func TestWhenNestedRulePointerIsMutatedOnTheClone_SourceRuleIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state.ZoneContentRow{
		Sid:   gofakeit.LetterN(10),
		Count: 1,
		Rules: []editor_state.ContentRuleRow{{Name: gofakeit.LetterN(8), IsGuarded: new(true)}},
	}
	clone := editor_state_helpers.CloneZoneContentRow(row)

	// Act
	*clone.Rules[0].IsGuarded = false

	// Assert
	require.NotNil(t, row.Rules[0].IsGuarded)
	assert.True(t, *row.Rules[0].IsGuarded)
}
