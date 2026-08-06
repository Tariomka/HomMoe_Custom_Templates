package zoneEditorHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDeletionIsChecked_ReturnsTheEditorsVerdict(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	zoneName := gofakeit.Word()
	playerZoneNames := map[string]bool{gofakeit.Word(): true}
	fixture.zoneEditor.On("CanDeleteZone", zoneName, playerZoneNames).Return(false)

	// Act
	canDelete := fixture.handler.CanDeleteZone(zoneName, playerZoneNames)

	// Assert
	assert.False(t, canDelete)
}
