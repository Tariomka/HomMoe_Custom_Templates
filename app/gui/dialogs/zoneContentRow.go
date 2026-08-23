package dialogs

import (
	"slices"

	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// zoneContentRow is one editable item inside a zone-content section.
type zoneContentRow struct {
	Mapping models.SidMapping
	Count   int
	IsGroup bool
	rules   []models.ContentRuleRow

	countSld  widget.Float
	manageBtn widget.Clickable
	removeBtn widget.Clickable
	dupBtn    widget.Clickable
}

func newZoneContentRow(
	mapping models.SidMapping,
	count int,
	rules []models.ContentRuleRow,
	isGroup bool) *zoneContentRow {
	return &zoneContentRow{
		Mapping: mapping,
		Count:   count,
		IsGroup: isGroup,
		rules:   slices.Clone(rules),
	}
}
