package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenResourcesContentPoolIsProvided_SetsResourcesContentPoolOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedPool := []string{gofakeit.Word(), gofakeit.Word()}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithResourcesContentPool(expectedPool).Build()

	// Assert
	assert.Equal(t, entities.Zone{ResourcesContentPool: expectedPool}, zone)
}
