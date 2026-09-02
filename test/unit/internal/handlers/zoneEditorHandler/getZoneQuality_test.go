package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneQualityIsRequested_ReturnsTheClassifiersQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	zone := entities.Zone{Name: gofakeit.Word()}
	fixture.tierService.On("GetQuality", zone).Return(neutral_zone.QualityHigh)

	// Act
	quality := fixture.handler.GetZoneQuality(zone)

	// Assert
	assert.Equal(t, neutral_zone.QualityHigh, quality)
}
