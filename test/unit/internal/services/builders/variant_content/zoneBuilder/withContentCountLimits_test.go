package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenContentCountLimitsAreProvided_SetsContentCountLimitsOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedLimits := []string{gofakeit.Word(), gofakeit.Word()}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithContentCountLimits(expectedLimits).Build()

	// Assert
	assert.Equal(t, entities.Zone{ContentCountLimits: entities.StringList(expectedLimits)}, zone)
}
