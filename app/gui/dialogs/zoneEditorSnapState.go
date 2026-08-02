package dialogs

import "gioui.org/widget"

type zoneEditorSnapState struct {
	snapBool         widget.Bool
	snapGuideX       float64
	snapGuideY       float64
	snapGuideXActive bool
	snapGuideYActive bool
}
