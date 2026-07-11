package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenCenterLayoutIsChosen_SetsCenterLayoutOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithLayoutCenter().Build()

	// Assert
	assert.Equal(t, entities.Zone{Layout: "zone_layout_center"}, zone)
}
