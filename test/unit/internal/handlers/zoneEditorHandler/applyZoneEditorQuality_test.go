package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenQualityIsApplied_ReturnsTheMutatedZone(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	appliedName := gofakeit.Word()
	fixture.zoneEditor.
		On("ApplyNeutralZoneQuality", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(arguments mock.Arguments) {
			zone, _ := arguments.Get(0).(*entities.Zone)
			zone.Name = appliedName
		}).
		Return()

	// Act
	zone := fixture.handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone: entities.Zone{Name: gofakeit.Word()},
	})

	// Assert
	assert.Equal(t, appliedName, zone.Name)
}

func TestWhenQualityIsApplied_ForwardsTheRequestedQualityAndCastleCount(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	castleCount := gofakeit.IntRange(0, 4)
	tuning := models.GenerationTuning{ContentScale: gofakeit.Float64Range(0.5, 2)}
	fixture.zoneEditor.
		On("ApplyNeutralZoneQuality", mock.Anything, neutral_zone.QualityLow, castleCount, tuning).
		Return()

	// Act
	_ = fixture.handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone:        entities.Zone{Name: gofakeit.Word()},
		Quality:     neutral_zone.QualityLow,
		CastleCount: castleCount,
		Tuning:      tuning,
	})

	// Assert
	fixture.zoneEditor.AssertCalled(
		t, "ApplyNeutralZoneQuality", mock.Anything, neutral_zone.QualityLow, castleCount, tuning)
}
