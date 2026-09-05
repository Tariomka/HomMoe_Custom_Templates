package dialogs

import (
	"slices"

	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// zoneContentRow is one editable item inside a zone-content section.
type zoneContentRow struct {
	Mapping models.SidMapping
	Count   int
	IsGroup bool
	rules   []editor_state_model.ContentRuleRow

	countSld  widget.Float
	manageBtn widget.Clickable
	removeBtn widget.Clickable
	dupBtn    widget.Clickable
}

func newZoneContentRow(
	mapping models.SidMapping,
	count int,
	rules []editor_state_model.ContentRuleRow,
	isGroup bool) *zoneContentRow {
	return &zoneContentRow{
		Mapping: mapping,
		Count:   count,
		IsGroup: isGroup,
		rules:   slices.Clone(rules),
	}
}
