package zoneTierService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheZoneCarriesATier_ReturnsTheRecordedTier(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zone := template_model.Zone{
		Name:               "Neutral-A",
		Quality:            new(neutral_zone.QualityHighest),
		Layout:             registry.GetLayoutValues().Sides,
		GuardedContentPool: []string{"pool_t1_stuff"},
	}

	// Act
	quality := service.ResolveQuality(zone)

	// Assert
	assert.Equal(t, neutral_zone.QualityHighest, quality)
}

func TestWhenTheZoneCarriesTheLowestTier_ReturnsItInsteadOfInferring(t *testing.T) {
	t.Parallel()
	// Arrange
	// A value field would make this case indistinguishable from "nothing
	// recorded", because the Quality enum counts from iota - 1.
	service := zones.NewZoneTierService()
	zone := template_model.Zone{
		Name:               "Neutral-A",
		Quality:            new(neutral_zone.QualityLowest),
		Layout:             registry.GetLayoutValues().TreasureZone,
		GuardedContentPool: []string{"pool_t3_stuff"},
	}

	// Act
	quality := service.ResolveQuality(zone)

	// Assert
	assert.Equal(t, neutral_zone.QualityLowest, quality)
}

func TestWhenTheZoneCarriesNoTier_InfersTheTier(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zone := template_model.Zone{
		Name:               "Neutral-A",
		Layout:             registry.GetLayoutValues().TreasureZone,
		GuardedContentPool: []string{"pool_t3_stuff"},
	}

	// Act
	quality := service.ResolveQuality(zone)

	// Assert
	assert.Equal(t, neutral_zone.QualityMedium, quality)
}

func TestWhenTheZoneCarriesNoTierAndCannotBeInferred_ReturnsUnknown(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zone := template_model.Zone{Name: "Neutral-A"}

	// Act
	quality := service.ResolveQuality(zone)

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}
