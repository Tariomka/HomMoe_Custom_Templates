package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenEncounterHolesSettingsAreProvided_SetsEncounterHolesSettingsPointerOnBuiltZone(t *testing.T) {
	// Arrange
	expectedSettings := entities.EncounterHolesSettings{
		AffectedEncounters: gofakeit.Float64Range(0.01, 1),
		TwoHoleEncounters:  gofakeit.Float64Range(0.01, 1),
	}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithEncounterHolesSettings(expectedSettings).Build()

	// Assert
	assert.Equal(t, entities.Zone{EncounterHolesSettings: &expectedSettings}, zone)
}
