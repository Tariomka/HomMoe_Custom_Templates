package variantMapping_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenFormattedAsString_ReturnsDisplayText(t *testing.T) {
	t.Parallel()
	// Arrange
	mapping := models.NewVariantMapping(
		models.SidMapping{Sid: "x", Name: "Fallback"},
		map[int]string{2: "second", 0: "first", 1: "middle"},
	)

	// Act
	stringValue := mapping.String()

	// Assert
	assert.Equal(t, "first", stringValue)
}
