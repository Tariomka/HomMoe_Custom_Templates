package dialogs

import (
	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
)

type zoneEditorZonePropertiesState struct {
	syncedZoneFor   string
	qualityDropdown *components.DropdownSelector
	castleDropdown  *components.DropdownSelector
	zoneSizeEdit    widget.Editor
	zoneGuardEdit   widget.Editor
	zoneWeeklyEdit  widget.Editor
	sideZoneDelete  widget.Clickable
}
