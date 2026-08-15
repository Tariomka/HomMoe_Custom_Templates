package zoneContentRowSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRowIsCloned_ScalarFieldsAreCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	row := models.ZoneContentRowSave{Sid: gofakeit.LetterN(10), Count: gofakeit.IntRange(1, 9), IsGroup: true}

	// Act
	clone := row.Clone()

	// Assert
	assert.Equal(t, row, clone)
}

func TestWhenRulesAreNil_CloneRulesStayNil(t *testing.T) {
	t.Parallel()
	// Arrange - a nil slice must not become an empty one, because the editor's
	// change detection distinguishes the two.
	row := models.ZoneContentRowSave{Sid: gofakeit.LetterN(10), Count: 1}

	// Act
	clone := row.Clone()

	// Assert
	assert.Nil(t, clone.Rules)
}

func TestWhenRuleIsReplacedInPlaceOnTheClone_SourceRuleIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	ruleName := gofakeit.LetterN(8)
	row := models.ZoneContentRowSave{
		Sid:   gofakeit.LetterN(10),
		Count: 1,
		Rules: []models.ContentRuleRowSave{{Name: ruleName}},
	}
	clone := row.Clone()

	// Act
	clone.Rules[0].Name = ruleName + "-changed"

	// Assert
	assert.Equal(t, ruleName, row.Rules[0].Name)
}

func TestWhenNestedRulePointerIsMutatedOnTheClone_SourceRuleIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	row := models.ZoneContentRowSave{
		Sid:   gofakeit.LetterN(10),
		Count: 1,
		Rules: []models.ContentRuleRowSave{{Name: gofakeit.LetterN(8), IsGuarded: new(true)}},
	}
	clone := row.Clone()

	// Act
	*clone.Rules[0].IsGuarded = false

	// Assert
	require.NotNil(t, row.Rules[0].IsGuarded)
	assert.True(t, *row.Rules[0].IsGuarded)
}
