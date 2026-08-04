package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactoriesAreProvided_ReturnsUsableService(t *testing.T) {
	t.Parallel()
	// Arrange
	castleFactory := zones.NewCastleFactory()
	roadFactory := zones.NewRoadFactory()

	// Act
	service := connection_editor.NewZoneEditorService(
		castleFactory,
		roadFactory,
		zones.NewZoneFactory(castleFactory, roadFactory),
	)

	// Assert
	assert.NotNil(t, service)
}
