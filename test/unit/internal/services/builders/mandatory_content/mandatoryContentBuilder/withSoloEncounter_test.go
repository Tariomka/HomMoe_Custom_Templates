package mandatoryContentBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSoloEncounterIsChosen_MarksBuiltItemAsSoloEncounter(t *testing.T) {
	// Arrange
	expectedSid := gofakeit.Word()
	builder := mandatory_content.NewContentBuilder(expectedSid)

	// Act
	item := builder.WithSoloEncounter().Build()

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{SID: expectedSid, SoloEncounter: true}, item)
}
