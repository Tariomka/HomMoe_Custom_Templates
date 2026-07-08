package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

// WeeklyIncrementValues is a public preset table (not a function); it is pinned
// here because NewDefaultConnection and the GUI dropdown both depend on it.
func TestWhenReadingWeeklyIncrementPresets_MatchesGuardGrowthTable(t *testing.T) {
	// Arrange
	expected := []float64{0.05, 0.10, 0.15, 0.20, 0.25}

	// Act
	values := connection_editor.WeeklyIncrementValues

	// Assert
	assert.Equal(t, expected, values)
}
