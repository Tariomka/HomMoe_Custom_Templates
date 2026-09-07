package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsCreated_ReturnsTheEditorsZone(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorNeutralZoneRequestDto{
		Label:         gofakeit.Word(),
		Quality:       neutral_zone.QualityHigh,
		CastleCount:   gofakeit.IntRange(0, 4),
		GenerateRoads: true,
		Tuning:        models.GenerationTuning{ContentScale: gofakeit.Float64Range(0.5, 2)},
	}
	expected := template_model.Zone{Name: request.Label}
	fixture.zoneEditor.
		On("NewDefaultNeutralZone",
			request.Label, request.Quality, request.CastleCount, request.GenerateRoads, request.Tuning).
		Return(expected)

	// Act
	zone := fixture.handler.CreateZoneEditorNeutralZone(request)

	// Assert
	assert.Equal(t, expected, zone)
}
