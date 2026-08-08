package bonusHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSpellCountLabelIsRequested_ReturnsTheServiceCaption(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newBonusHandlerFixture()
	count := gofakeit.Number(0, 20)
	expected := gofakeit.Sentence(3)
	fixture.bonusService.On("GetSpellCountLabel", count).Return(expected)

	// Act
	label := fixture.handler.GetSpellCountLabel(count)

	// Assert
	assert.Equal(t, expected, label)
}
