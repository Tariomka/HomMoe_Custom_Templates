package components

import (
	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/content"
)

type BasicSetupPanel struct {
	templateName widget.Editor
	gameMode     *content.SegmentButtonGroup
	playerCnt    widget.Float
	mapSizeSld   widget.Float
	chkExpSizes  widget.Bool
	topology     *content.DropdownSelector
}
