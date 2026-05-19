package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type ZoneContentPanel struct {
	zcMines      *content.ZoneContentSection
	zcTreasures  *content.ZoneContentSection
	zcHires      *content.ZoneContentSection
	zcBanks      *content.ZoneContentSection
	btnZoneReset widget.Clickable

	scroll widget.List

	state *State
}

func NewZoneContentPanel(state *State) *ZoneContentPanel {
	panel := &ZoneContentPanel{

		zcMines:     content.NewZoneContentSection("Mines", constants.ContentItemGroup.Mines, 3, true),
		zcTreasures: content.NewZoneContentSection("Treasures", constants.ContentItemGroup.Treasures, 10, false),
		zcHires:     content.NewZoneContentSection("Random Hires", constants.ContentItemGroup.HireBuildings, 10, false),
		zcBanks:     content.NewZoneContentSection("Resource Banks", constants.ContentItemGroup.ResourceBanks, 10, false),
		state:       state,
	}
	panel.scroll.Axis = layout.Vertical
	panel.seedDefaultPlayerZoneContent()
	panel.LoadFromState()
	return panel
}

func (this *ZoneContentPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{}
	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

func (this *ZoneContentPanel) LoadFromState() {
	settings := this.state.GetSettingsFile()
	if len(settings.PlayerZoneMandatoryContent) > 0 {
		this.applyZoneContentItems(settings.PlayerZoneMandatoryContent)
	}
}

// TODO: check `.Update(gtx)` and on true update the value
func (this *ZoneContentPanel) SaveToState() {
	this.state.UpdateState(func(settings *models.SettingsFile) {
		settings.PlayerZoneMandatoryContent = this.collectZoneContentItems()
	})
}

// seedDefaultPlayerZoneContent mirrors C# InitializeDefaultPlayerZoneContents.
func (this *ZoneContentPanel) seedDefaultPlayerZoneContent() {
	this.zcMines.ClearRows()
	this.zcTreasures.ClearRows()
	this.zcHires.ClearRows()
	this.zcBanks.ClearRows()

	// Mines: wood/ore/gold guarded next-to-castle; crystals/mercury/gemstones/alchemy-lab guarded near road.
	this.zcMines.Add(constants.ContentIds.MineWood, 1, true, true, 0, false)
	this.zcMines.Add(constants.ContentIds.MineOre, 1, true, true, 0, false)
	this.zcMines.Add(constants.ContentIds.MineGold, 1, true, true, 0, false)
	this.zcMines.Add(constants.ContentIds.MineCrystals, 1, true, false, 1, false)
	this.zcMines.Add(constants.ContentIds.MineMercury, 1, true, false, 1, false)
	this.zcMines.Add(constants.ContentIds.MineGemstones, 1, true, false, 1, false)
	this.zcMines.Add(constants.ContentIds.AlchemyLab, 1, true, false, 1, false)

	// Treasures: PandoraBox + RandomItemEpic guarded.
	this.zcTreasures.Add(constants.ContentIds.PandoraBox, 1, true, false, 0, false)
	this.zcTreasures.Add(constants.ContentIds.RandomItemEpic, 1, true, false, 0, false)

	// Random hires: low ×2, high ×1, all-tier ×1 (groups).
	this.zcHires.Add(constants.IncludeListIds.RandomHiresLowTier, 2, true, false, 0, true)
	this.zcHires.Add(constants.IncludeListIds.RandomHiresHighTier, 1, true, false, 0, true)
	this.zcHires.Add(constants.IncludeListIds.RandomHiresAllTier, 1, true, false, 0, true)

	// Resource banks: tier1 ×2, tier2 ×1.
	this.zcBanks.Add(constants.IncludeListIds.ResourceBanksTier1, 2, true, false, 0, true)
	this.zcBanks.Add(constants.IncludeListIds.ResourceBanksTier2, 1, true, false, 0, true)
}

// applyZoneContentItems replaces every section based on a flat list of items
// loaded from a settings file. Items are routed to the appropriate section by
// SID lookup.
func (this *ZoneContentPanel) applyZoneContentItems(items []models.ZoneContentItem) {
	this.zcMines.ClearRows()
	this.zcTreasures.ClearRows()
	this.zcHires.ClearRows()
	this.zcBanks.ClearRows()
	for _, item := range items {
		mapping := models.SidMapping{Sid: item.Sid, Name: item.Name}
		if found, ok := helpers.LookupSid(item.Sid); ok {
			mapping = found
		}
		count := max(item.Count, 1)
		roadIdx := max(indexOf(constants.RoadDistances, item.RoadDistance), 0)
		switch {
		case sectionContains(constants.ContentItemGroup.Mines, item.Sid):
			this.zcMines.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		case sectionContains(constants.ContentItemGroup.Treasures, item.Sid):
			this.zcTreasures.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		case sectionContains(constants.ContentItemGroup.HireBuildings, item.Sid):
			this.zcHires.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		case sectionContains(constants.ContentItemGroup.ResourceBanks, item.Sid):
			this.zcBanks.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		default:
			// Unknown SID — keep with treasures by default.
			this.zcTreasures.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		}
	}
}

// collectZoneContentItems serialises every section into a flat list.
func (this *ZoneContentPanel) collectZoneContentItems() []models.ZoneContentItem {
	var out []models.ZoneContentItem
	collect := func(contentSection *content.ZoneContentSection) {
		for row := range contentSection.IterateRows() {
			roadDistance := ""
			if row.RoadDistIdx >= 0 && row.RoadDistIdx < len(constants.RoadDistances) {
				roadDistance = constants.RoadDistances[row.RoadDistIdx]
			}
			out = append(out, models.ZoneContentItem{
				Sid:          row.Mapping.Sid,
				Name:         row.Mapping.Name,
				Count:        row.Count,
				IsGuarded:    row.IsGuarded.Value,
				NearCastle:   row.NearCastle.Value,
				RoadDistance: roadDistance,
				IsGroup:      row.IsGroup,
			})
		}
	}
	collect(this.zcMines)
	collect(this.zcTreasures)
	collect(this.zcHires)
	collect(this.zcBanks)
	return out
}

// indexOf returns the index of value in items, or -1 when not present.
func indexOf[T comparable](items []T, value T) int {
	for i, candidate := range items {
		if candidate == value {
			return i
		}
	}
	return -1
}

func sectionContains(list []models.SidMapping, sid string) bool {
	for _, mapping := range list {
		if mapping.Sid == sid {
			return true
		}
	}
	return false
}
