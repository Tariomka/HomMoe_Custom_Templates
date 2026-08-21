package dialogs

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// ZoneContentDialog is a single-tier version of the former Zone Content tab. It
// edits the mandatory-content rows for one zone type and applies the result
// live through the onApply callback, mirroring the tab's behaviour.
type ZoneContentDialog struct {
	zcMines           *ZoneContentSection
	zcUtilities       *ZoneContentSection
	zcTreasures       *ZoneContentSection
	zcHires           *ZoneContentSection
	zcBanks           *ZoneContentSection
	zcHeroImprovement *ZoneContentSection

	title        string
	isPlayerTier bool
	onApply      func([]models.ZoneContentRowSave)

	btnReset widget.Clickable
	scroll   widget.List
}

// NewZoneContentDialog builds the dialog for a single zone tier. title is shown
// in the modal header, isPlayerTier seeds the reset defaults, rows are the
// current content, opener wires the per-row Manage Rules modal, and onApply is
// invoked with the edited rows whenever they change.
func NewZoneContentDialog(
	title string,
	isPlayerTier bool,
	rows []models.ZoneContentRowSave,
	contentRuleHandler handler_interfaces.IZoneContentHandler,
	opener interfaces.DialogOpener,
	onApply func([]models.ZoneContentRowSave),
) *ZoneContentDialog {
	dialog := &ZoneContentDialog{
		zcMines: NewZoneContentSection(
			"Mines", constants.ContentItemGroup.Mines, 3, true, contentRuleHandler),
		zcUtilities: NewZoneContentSection("Utility Structures",
			constants.ContentItemGroup.UtilityStructures, 10, false, contentRuleHandler),
		zcTreasures: NewZoneContentSection(
			"Treasures", constants.ContentItemGroup.Treasures, 10, false, contentRuleHandler),
		zcHires: NewZoneContentSection("Unit Recruitment",
			constants.ContentItemGroup.UnitRecruitment, 10, false, contentRuleHandler),
		zcBanks: NewZoneContentSection(
			"Resource Banks", constants.ContentItemGroup.ResourceBanks, 10, false, contentRuleHandler),
		zcHeroImprovement: NewZoneContentSection("Hero Improvement",
			constants.ContentItemGroup.HeroImprovementStructures, 10, false, contentRuleHandler),
		title:        title,
		isPlayerTier: isPlayerTier,
		onApply:      onApply,
	}
	dialog.scroll.Axis = layout.Vertical

	dialog.zcMines.SetDialogOpener(opener)
	dialog.zcUtilities.SetDialogOpener(opener)
	dialog.zcTreasures.SetDialogOpener(opener)
	dialog.zcHires.SetDialogOpener(opener)
	dialog.zcBanks.SetDialogOpener(opener)
	dialog.zcHeroImprovement.SetDialogOpener(opener)

	dialog.loadRowsIntoSections(rows)
	return dialog
}

func (this *ZoneContentDialog) Title() string {
	return this.title
}

func (this *ZoneContentDialog) PreferredSize() (unit.Dp, unit.Dp) {
	return unit.Dp(640), unit.Dp(560)
}

func (this *ZoneContentDialog) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	if this.btnReset.Clicked(gtx) {
		this.resetToDefault()
	}

	widgetsList := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(widgets.NewButtonWidget(theme, "Reset to default", &this.btnReset, false)),
			)
		},
		this.zcMines.Layout(theme),
		this.zcUtilities.Layout(theme),
		this.zcTreasures.Layout(theme),
		this.zcHires.Layout(theme),
		this.zcBanks.Layout(theme),
		this.zcHeroImprovement.Layout(theme),
	}

	dims := material.List(theme, &this.scroll).Layout(
		gtx, len(widgetsList),
		func(gtx layout.Context, index int) layout.Dimensions { return widgetsList[index](gtx) })

	// Persist after the sections have processed this frame's clicks so the
	// callback sees the latest edits, matching the live behaviour of the
	// former Zone Content tab.
	this.persist()
	return dims, false
}

// persist collects the current section rows and pushes them through onApply.
func (this *ZoneContentDialog) persist() {
	if this.onApply != nil {
		this.onApply(this.collectSectionRows())
	}
}

// resetToDefault reverts the sections to the tier defaults - the historical
// seeded defaults for the Player tier, otherwise an empty list.
func (this *ZoneContentDialog) resetToDefault() {
	if this.isPlayerTier {
		this.loadRowsIntoSections(editor_state_dto.DefaultPlayerZoneContentRows())
		return
	}
	this.loadRowsIntoSections(nil)
}

// loadRowsIntoSections replaces the section rows with the given list, routing
// each row to its appropriate section.
func (this *ZoneContentDialog) loadRowsIntoSections(rows []models.ZoneContentRowSave) {
	this.zcMines.ClearRows()
	this.zcUtilities.ClearRows()
	this.zcTreasures.ClearRows()
	this.zcHires.ClearRows()
	this.zcBanks.ClearRows()
	this.zcHeroImprovement.ClearRows()
	for _, raw := range rows {
		row := raw.Normalized()
		mapping := models.SidMapping{Sid: row.Sid, Name: row.Sid}
		if found, ok := utils.GetSidMappingBySid(row.Sid); ok {
			mapping = found
		}
		var rules []models.ContentRuleRowSave
		if len(row.Rules) > 0 {
			rules = row.Rules
		}
		section := this.routeToSection(row.Sid, row.IsMine)
		section.Add(mapping, row.Count, rules, row.IsGroup)
	}
}

// routeToSection picks the correct UI section for a SID, honouring the IsMine
// flag so that a mine-include-list still ends up in Mines.
func (this *ZoneContentDialog) routeToSection(sid string, isMine bool) *ZoneContentSection {
	switch {
	case isMine || sectionContains(constants.ContentItemGroup.Mines, sid):
		return this.zcMines
	case sectionContains(constants.ContentItemGroup.UnitRecruitment, sid):
		return this.zcHires
	case sectionContains(constants.ContentItemGroup.ResourceBanks, sid):
		return this.zcBanks
	case sectionContains(constants.ContentItemGroup.HeroImprovementStructures, sid):
		return this.zcHeroImprovement
	case sectionContains(constants.ContentItemGroup.UtilityStructures, sid):
		return this.zcUtilities
	case sectionContains(constants.ContentItemGroup.Treasures, sid):
		return this.zcTreasures
	default:
		return this.zcTreasures
	}
}

// collectSectionRows reads the current sections back into a flat slice of
// save-rows tagged with the correct IsMine flag.
func (this *ZoneContentDialog) collectSectionRows() []models.ZoneContentRowSave {
	var out []models.ZoneContentRowSave
	gather := func(section *ZoneContentSection, isMine bool) {
		for row := range section.IterateRows() {
			out = append(out, models.ZoneContentRowSave{
				Sid:     row.Mapping.Sid,
				Count:   row.Count,
				IsGroup: row.IsGroup,
				IsMine:  isMine,
				Rules:   row.Rules(),
			})
		}
	}
	gather(this.zcMines, true)
	gather(this.zcUtilities, false)
	gather(this.zcTreasures, false)
	gather(this.zcHires, false)
	gather(this.zcBanks, false)
	gather(this.zcHeroImprovement, false)
	return out
}

func sectionContains(list []models.SidMapping, sid string) bool {
	for _, mapping := range list {
		if mapping.Sid == sid {
			return true
		}
	}
	return false
}
