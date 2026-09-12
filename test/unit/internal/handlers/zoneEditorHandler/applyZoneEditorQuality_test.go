package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenQualityIsApplied_ReturnsTheMutatedZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	mutated := []template_model.Zone{{Name: gofakeit.Word()}}
	fixture.zoneEditor.
		On("ApplyNeutralZoneQualityEdit", mock.Anything).
		Return(mutated, []template_model.Connection(nil))

	// Act
	mutation := fixture.handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone: template_model.Zone{Name: gofakeit.Word()},
	})

	// Assert
	assert.Equal(t, mutated, mutation.Zones)
}

func TestWhenQualityIsApplied_ReturnsTheMutatedConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	mutated := []template_model.Connection{{Name: gofakeit.Word(), GuardValue: gofakeit.IntRange(1, 60000)}}
	fixture.zoneEditor.
		On("ApplyNeutralZoneQualityEdit", mock.Anything).
		Return([]template_model.Zone(nil), mutated)

	// Act
	mutation := fixture.handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone: template_model.Zone{Name: gofakeit.Word()},
	})

	// Assert
	assert.Equal(t, mutated, mutation.Connections)
}

func TestWhenQualityIsApplied_ForwardsTheWorkingGraphAndTheRequestedTier(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	zoneName := gofakeit.Word()
	expected := models.NeutralZoneQualityEditRequest{
		Zones:           []template_model.Zone{{Name: zoneName}},
		Connections:     []template_model.Connection{{Name: gofakeit.Word(), From: zoneName}},
		PlayerZoneNames: []string{gofakeit.Word()},
		ZoneName:        zoneName,
		Quality:         neutral_zone.QualityLow,
		CastleCount:     gofakeit.IntRange(0, 4),
		Tuning:          models.GenerationTuning{ContentScale: gofakeit.Float64Range(0.5, 2)},
	}
	fixture.zoneEditor.
		On("ApplyNeutralZoneQualityEdit", expected).
		Return([]template_model.Zone(nil), []template_model.Connection(nil))

	// Act
	_ = fixture.handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone:            template_model.Zone{Name: zoneName},
		Quality:         expected.Quality,
		CastleCount:     expected.CastleCount,
		Tuning:          expected.Tuning,
		Zones:           expected.Zones,
		Connections:     expected.Connections,
		PlayerZoneNames: expected.PlayerZoneNames,
	})

	// Assert
	fixture.zoneEditor.AssertCalled(t, "ApplyNeutralZoneQualityEdit", expected)
}
