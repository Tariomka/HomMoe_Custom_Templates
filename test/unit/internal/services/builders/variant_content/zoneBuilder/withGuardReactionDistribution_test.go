package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardReactionDistributionIsProvided_SetsGuardReactionDistributionOnBuiltZone(t *testing.T) {
	// Arrange
	expectedDistribution := []int{gofakeit.Number(1, 100), gofakeit.Number(1, 100)}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithGuardReactionDistribution(expectedDistribution).Build()

	// Assert
	assert.Equal(t, entities.Zone{GuardReactionDistribution: expectedDistribution}, zone)
}
