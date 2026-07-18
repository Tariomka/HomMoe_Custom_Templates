package mandatoryContentItemBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMineIsChosen_MarksBuiltItemAsMine(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSid := gofakeit.Word()
	builder := mandatory_content.NewContentItemBuilder(expectedSid)

	// Act
	item := builder.WithMine().Build()

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{SID: expectedSid, IsMine: true}, item)
}
