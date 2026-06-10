package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
)

// tierIndex is the position of a tier inside the tierRows cache. The
// numeric order must match the SegmentButtonGroup label list below
type tierIndex int

const (
	tierPlayer tierIndex = iota
	tierLow
	tierMedium
	tierHigh
	tierHub
	tierCount
)

var tierLabels = []string{"Player", "Low Neutral", "Medium Neutral", "High Neutral", "Hub"}

type ZoneContentPanel struct {
	zcMines           *content.ZoneContentSection
	zcUtilities       *content.ZoneContentSection
	zcTreasures       *content.ZoneContentSection
	zcHires           *content.ZoneContentSection
	zcBanks           *content.ZoneContentSection
	zcHeroImprovement *content.ZoneContentSection

	tierSelector *content.SegmentButtonGroup
	currentTier  tierIndex
	// tierRows caches the rows for every tier so the user can swap
	// between tabs without losing in-progress edits.
	tierRows [tierCount][]models.ZoneContentRowSave

	btnZoneReset widget.Clickable

	scroll widget.List

	state *State
}

func NewZoneContentPanel(state *State) *ZoneContentPanel {
	panel := &ZoneContentPanel{
		zcMines:           content.NewZoneContentSection("Mines", constants.ContentItemGroup.Mines, 3, true),
		zcUtilities:       content.NewZoneContentSection("Utility Structures", constants.ContentItemGroup.UtilityStructures, 10, false),
		zcTreasures:       content.NewZoneContentSection("Treasures", constants.ContentItemGroup.Treasures, 10, false),
		zcHires:           content.NewZoneContentSection("Unit Recruitment", constants.ContentItemGroup.UnitRecruitment, 10, false),
		zcBanks:           content.NewZoneContentSection("Resource Banks", constants.ContentItemGroup.ResourceBanks, 10, false),
		zcHeroImprovement: content.NewZoneContentSection("Hero Improvement", constants.ContentItemGroup.HeroImprovementStructures, 10, false),
		tierSelector:      content.NewSegmentButtonGroup(tierLabels),
		state:             state,
	}
	panel.scroll.Axis = layout.Vertical

	// Wire every section to the modal host so rows can open the Manage Rules dialog.
	opener := state.Dialogs().Open
	panel.zcMines.SetDialogOpener(opener)
	panel.zcUtilities.SetDialogOpener(opener)
	panel.zcTreasures.SetDialogOpener(opener)
	panel.zcHires.SetDialogOpener(opener)
	panel.zcBanks.SetDialogOpener(opener)
	panel.zcHeroImprovement.SetDialogOpener(opener)

	panel.LoadFromState()
	return panel
}

func (this *ZoneContentPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Reset button is global to the currently-selected tier.
		if this.btnZoneReset.Clicked(gtx) {
			this.resetCurrentTier()
		}
		// Tier switch: persist current section into cache, then load the
		// new tier's rows into the shared sections.
		if this.tierSelector.Update(gtx) {
			this.cacheCurrentSections()
			this.currentTier = tierIndex(this.tierSelector.GetSelectedIndex())
			this.loadTierIntoSections(this.currentTier)
		}

		widgetsList := []layout.Widget{
			widgets.NewWarningBannerWidget(theme, "EXPERIMENTAL — Mandatory content per zone tier. Effects only apply on generation."),
			func(gtx layout.Context) layout.Dimensions {
				return this.tierSelector.Layout(gtx, theme)
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(widgets.NewButtonWidget(theme, "↺  Reset this tier", &this.btnZoneReset, false)),
					)
				})
			},
			this.zcMines.Layout(theme),
			this.zcUtilities.Layout(theme),
			this.zcTreasures.Layout(theme),
			this.zcHires.Layout(theme),
			this.zcBanks.Layout(theme),
			this.zcHeroImprovement.Layout(theme),
		}
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

// LoadFromState pulls all five tier lists out of the SettingsFile,
// caches them, and shows the currently-selected tier.
func (this *ZoneContentPanel) LoadFromState() {
	settings := this.state.GetStateData()
	this.tierRows[tierPlayer] = append([]models.ZoneContentRowSave(nil), settings.PlayerZoneContentRows...)
	this.tierRows[tierLow] = append([]models.ZoneContentRowSave(nil), settings.LowNeutralContentRows...)
	this.tierRows[tierMedium] = append([]models.ZoneContentRowSave(nil), settings.MediumNeutralContentRows...)
	this.tierRows[tierHigh] = append([]models.ZoneContentRowSave(nil), settings.HighNeutralContentRows...)
	this.tierRows[tierHub] = append([]models.ZoneContentRowSave(nil), settings.HubZoneContentRows...)

	if len(this.tierRows[tierPlayer]) == 0 {
		this.tierRows[tierPlayer] = defaultPlayerTierRows()
	}

	this.tierSelector.SetSelectedIndex(int(this.currentTier))
	this.loadTierIntoSections(this.currentTier)
}

// SaveToState collects every tier's rows back into the SettingsFile.
func (this *ZoneContentPanel) SaveToState() {
	this.cacheCurrentSections()
	this.state.UpdateState(func(settings *models.EditorStateModel) {
		settings.PlayerZoneContentRows = cloneRows(this.tierRows[tierPlayer])
		settings.LowNeutralContentRows = cloneRows(this.tierRows[tierLow])
		settings.MediumNeutralContentRows = cloneRows(this.tierRows[tierMedium])
		settings.HighNeutralContentRows = cloneRows(this.tierRows[tierHigh])
		settings.HubZoneContentRows = cloneRows(this.tierRows[tierHub])
	})
}

// loadTierIntoSections replaces the section rows with the given tier's
// cached row list, routing each row to its appropriate section.
func (this *ZoneContentPanel) loadTierIntoSections(tier tierIndex) {
	this.zcMines.ClearRows()
	this.zcUtilities.ClearRows()
	this.zcTreasures.ClearRows()
	this.zcHires.ClearRows()
	this.zcBanks.ClearRows()
	this.zcHeroImprovement.ClearRows()
	for _, raw := range this.tierRows[tier] {
		row := raw.Normalised()
		mapping := models.SidMapping{Sid: row.Sid, Name: row.Sid}
		if found, ok := utils.LookupSid(row.Sid); ok {
			mapping = found
		}
		// Prefer the explicit rule list; otherwise migrate legacy flat fields,
		// but only when they actually carry rule data so empty rows stay clean.
		var rules []models.ContentRuleRowSave
		switch {
		case len(row.Rules) > 0:
			rules = row.Rules
		case row.HasLegacyRuleData():
			rules = content_rules.MigrateLegacyRow(row, mapping).Rules
		}
		section := this.routeToSection(row.Sid, row.IsMine)
		section.Add(mapping, row.Count, rules, row.IsGroup)
	}
}

// cacheCurrentSections serializes the shared sections back into the
// row cache for the currently-active tier.
func (this *ZoneContentPanel) cacheCurrentSections() {
	this.tierRows[this.currentTier] = this.collectSectionRows()
}

// resetCurrentTier reverts the active tier to its defaults — that means
// the historical seeded defaults for Player, and an empty list for the
// other four tiers
func (this *ZoneContentPanel) resetCurrentTier() {
	switch this.currentTier {
	case tierPlayer:
		this.tierRows[tierPlayer] = defaultPlayerTierRows()
	default:
		this.tierRows[this.currentTier] = nil
	}
	this.loadTierIntoSections(this.currentTier)
}

// routeToSection picks the correct UI section for a SID, honouring the
// IsMine flag so that a mine-include-list still ends up in Mines.
func (this *ZoneContentPanel) routeToSection(sid string, isMine bool) *content.ZoneContentSection {
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

// collectSectionRows reads the current sections back into a flat slice
// of save-rows tagged with the correct IsMine flag.
func (this *ZoneContentPanel) collectSectionRows() []models.ZoneContentRowSave {
	var out []models.ZoneContentRowSave
	gather := func(section *content.ZoneContentSection, isMine bool) {
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

func cloneRows(rows []models.ZoneContentRowSave) []models.ZoneContentRowSave {
	if len(rows) == 0 {
		return nil
	}
	return append([]models.ZoneContentRowSave(nil), rows...)
}

func defaultPlayerTierRows() []models.ZoneContentRowSave {
	return []models.ZoneContentRowSave{
		{Sid: constants.ContentIds.MineWood.Sid, Count: 1, IsGuarded: true, NearCastle: true, RoadDistance: "Any", IsMine: true},
		{Sid: constants.ContentIds.MineOre.Sid, Count: 1, IsGuarded: true, NearCastle: true, RoadDistance: "Any", IsMine: true},
		{Sid: constants.ContentIds.MineGold.Sid, Count: 1, IsGuarded: true, NearCastle: true, RoadDistance: "Any", IsMine: true},
		{Sid: constants.ContentIds.MineCrystals.Sid, Count: 1, IsGuarded: true, RoadDistance: "Next To", IsMine: true},
		{Sid: constants.ContentIds.MineMercury.Sid, Count: 1, IsGuarded: true, RoadDistance: "Next To", IsMine: true},
		{Sid: constants.ContentIds.MineGemstones.Sid, Count: 1, IsGuarded: true, RoadDistance: "Next To", IsMine: true},
		{Sid: constants.ContentIds.AlchemyLab.Sid, Count: 1, IsGuarded: true, RoadDistance: "Next To", IsMine: true},
		{Sid: constants.ContentIds.PandoraBox.Sid, Count: 1, IsGuarded: true, RoadDistance: "Any"},
		{Sid: constants.ContentIds.RandomItemEpic.Sid, Count: 1, IsGuarded: true, RoadDistance: "Any"},
		{Sid: constants.IncludeListIds.RandomHiresLowTier.Sid, Count: 2, IsGuarded: true, RoadDistance: "Any", IsGroup: true},
		{Sid: constants.IncludeListIds.RandomHiresHighTier.Sid, Count: 1, IsGuarded: true, RoadDistance: "Any", IsGroup: true},
		{Sid: constants.IncludeListIds.RandomHiresAllTier.Sid, Count: 1, IsGuarded: true, RoadDistance: "Any", IsGroup: true},
		{Sid: constants.IncludeListIds.ResourceBanksTier1.Sid, Count: 2, IsGuarded: true, RoadDistance: "Any", IsGroup: true},
		{Sid: constants.IncludeListIds.ResourceBanksTier2.Sid, Count: 1, IsGuarded: true, RoadDistance: "Any", IsGroup: true},
	}
}

func sectionContains(list []models.SidMapping, sid string) bool {
	for _, mapping := range list {
		if mapping.Sid == sid {
			return true
		}
	}
	return false
}
