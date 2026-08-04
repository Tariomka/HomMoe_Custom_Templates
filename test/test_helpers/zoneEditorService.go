package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// NewZoneEditorService builds a ZoneEditorService with the same collaborators the
// composition root wires, for tests that need a working editor rather than a mock.
func NewZoneEditorService() *connection_editor.ZoneEditorService {
	castleFactory := zones.NewCastleFactory()
	roadFactory := zones.NewRoadFactory()

	return connection_editor.NewZoneEditorService(
		castleFactory,
		roadFactory,
		zones.NewZoneFactory(castleFactory, roadFactory),
	)
}
