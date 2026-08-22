package dialogs

import (
	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// zoneContentRow is one editable item inside a zone-content section.
type zoneContentRow struct {
	Mapping models.SidMapping
	Count   int
	IsGroup bool
	rules   []models.ContentRuleRowSave

	countSld  widget.Float
	manageBtn widget.Clickable
	removeBtn widget.Clickable
	dupBtn    widget.Clickable
}

func newZoneContentRow(
	mapping models.SidMapping,
	count int,
	rules []models.ContentRuleRowSave,
	isGroup bool) *zoneContentRow {
	return &zoneContentRow{
		Mapping: mapping,
		Count:   count,
		IsGroup: isGroup,
		rules:   utils.CloneRuleRows(rules),
	}
}

// Rules returns a defensive copy of the row's content rules, letting the parent
// panel serialize them without exposing the row's mutable slice.
func (this *zoneContentRow) Rules() []models.ContentRuleRowSave {
	return utils.CloneRuleRows(this.rules)
}
