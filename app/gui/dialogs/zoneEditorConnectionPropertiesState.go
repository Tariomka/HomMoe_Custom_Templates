package dialogs

import (
	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type zoneEditorConnectionPropertiesState struct {
	syncedFor         *entities.Connection
	typeDropdown      *components.DropdownSelector
	guardZoneDropdown *components.DropdownSelector
	guardDropdown     *components.DropdownSelector
	guardPresetValues []int
	weeklyDropdown    *components.DropdownSelector
	guardValueEdit    widget.Editor
	weeklyEdit        widget.Editor
	matchGroupEdit    widget.Editor
	advancedBool      widget.Bool
	escapeBool        widget.Bool
	simSquadBool      widget.Bool
	sidePropDelete    widget.Clickable
}
