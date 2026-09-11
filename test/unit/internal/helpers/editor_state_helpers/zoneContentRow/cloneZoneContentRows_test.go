package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/editor_state_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRowsAreCloned_EveryRowIsCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state.ZoneContentRow{
		{Sid: gofakeit.LetterN(10), Count: gofakeit.IntRange(1, 9)},
		{Sid: gofakeit.LetterN(10), Count: gofakeit.IntRange(1, 9)},
	}

	// Act
	clones := editor_state_helpers.CloneZoneContentRows(rows)

	// Assert
	assert.Equal(t, rows, clones)
}

func TestWhenRowsAreCloned_TheNestedRulesAreNotShared(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state.ZoneContentRow{{
		Sid:   gofakeit.LetterN(10),
		Count: 1,
		Rules: []editor_state.ContentRuleRow{{Name: gofakeit.LetterN(8)}},
	}}

	// Act
	clones := editor_state_helpers.CloneZoneContentRows(rows)

	// Assert
	require.Len(t, clones, 1)
	require.Len(t, clones[0].Rules, 1)
	assert.NotSame(t, &rows[0].Rules[0], &clones[0].Rules[0])
}
