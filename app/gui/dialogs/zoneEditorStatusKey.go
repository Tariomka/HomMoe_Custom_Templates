package dialogs

// zoneEditorStatusKey is everything the zone editor's status line is drawn
// from. Comparing it across a frame reveals edits made after the line was laid out.
type zoneEditorStatusKey struct {
	hint            string
	addMode         bool
	addZoneMode     bool
	zoneCount       int
	connectionCount int
	graphDirty      bool
}
