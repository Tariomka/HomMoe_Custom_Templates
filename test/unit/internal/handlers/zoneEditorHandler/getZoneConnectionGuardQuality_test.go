package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionGuardQualityIsRequested_ReturnsTheClassifiersQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	from := gofakeit.Word()
	to := gofakeit.Word()
	zones := []template_model.Zone{{Name: from}, {Name: to}}
	playerName := gofakeit.Word()
	fixture.tierService.
		On("GetConnectionGuardQuality", from, to, zones, []string{playerName}).
		Return(neutral_zone.QualityMedium)

	// Act
	quality := fixture.handler.GetZoneConnectionGuardQuality(
		from, to, zones, map[string]bool{playerName: true})

	// Assert
	assert.Equal(t, neutral_zone.QualityMedium, quality)
}

func TestWhenThereAreNoPlayerZones_PassesAnEmptyPlayerNameList(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	from := gofakeit.Word()
	to := gofakeit.Word()
	fixture.tierService.
		On("GetConnectionGuardQuality", from, to, []template_model.Zone(nil), []string{}).
		Return(neutral_zone.QualityUnknown)

	// Act
	_ = fixture.handler.GetZoneConnectionGuardQuality(from, to, nil, map[string]bool{})

	// Assert
	fixture.tierService.AssertCalled(
		t, "GetConnectionGuardQuality", from, to, []template_model.Zone(nil), []string{})
}
