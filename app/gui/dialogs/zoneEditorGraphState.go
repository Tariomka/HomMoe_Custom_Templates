package dialogs

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

// zoneEditorGraphState caches the graph diagnostics the status line shows. Its
// dirty flag is set by every edit that replaces the zone or connection lists,
// and is kept apart from geometryDirty, which hit tests clear mid-input.
type zoneEditorGraphState struct {
	graph      dtos.ZoneEditorGraphDto
	graphDirty bool
	// shownStatus is the status key the status line was last drawn from.
	shownStatus zoneEditorStatusKey
}

func (this *zoneEditorGraphState) markGraphDirty() {
	this.graphDirty = true
}
