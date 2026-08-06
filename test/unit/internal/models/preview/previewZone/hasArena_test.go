package previewZone_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHostsTheArena_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := preview.Zone{Arena: true}

	// Act
	hasArena := zone.HasArena()

	// Assert
	assert.True(t, hasArena)
}

func TestWhenZoneDoesNotHostTheArena_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := preview.Zone{Arena: false}

	// Act
	hasArena := zone.HasArena()

	// Assert
	assert.False(t, hasArena)
}
