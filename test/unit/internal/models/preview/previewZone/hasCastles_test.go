package previewZone_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHasCastles_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := preview.Zone{Castles: gofakeit.Number(1, 4)}

	// Act
	hasCastles := zone.HasCastles()

	// Assert
	assert.True(t, hasCastles)
}

func TestWhenZoneHasNoCastles_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := preview.Zone{Castles: 0}

	// Act
	hasCastles := zone.HasCastles()

	// Assert
	assert.False(t, hasCastles)
}
