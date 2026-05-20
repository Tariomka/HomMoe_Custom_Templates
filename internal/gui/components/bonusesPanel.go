package components

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// BonusesPanel is the editor for the Game Start Bonuses, Banned Items,
// Banned Spells and Guard-Value Overrides lists added in C# v0.7.
//
// The C# WPF editor presents these as picker dialogs that build the
// underlying strings; the Go port exposes the raw newline-separated
// text directly so the user can copy/paste between editors and so we
// preserve full data parity without porting four modal dialogs.
type BonusesPanel struct {
	bonusesEdit       widget.Editor
	bannedItemsEdit   widget.Editor
	bannedMagicsEdit  widget.Editor
	valueOverrideEdit widget.Editor

	scroll widget.List

	state *State
}

func NewBonusesPanel(state *State) *BonusesPanel {
	panel := &BonusesPanel{state: state}
	panel.scroll.Axis = layout.Vertical
	panel.LoadFromState()
	return panel
}

func (this *BonusesPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{
		widgets.NewWarningBannerWidget(theme, "EXPERIMENTAL — Effects only apply on generation."),

		widgets.NewSectionWidget(theme, "Game start bonuses", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One bonus per line. Format: PresetType|ReceiverFilter|Param|Param2"),
			widgets.NewDimmedLabelWidget(theme, "PresetType: TownPortalFree, Spell, UnitMultiplier, MovementBonus, StartingItem, StartingGold, StartingGems, StartingCrystals, StartingMercury, StartingWood, StartingOre"),
			widgets.NewDimmedLabelWidget(theme, "ReceiverFilter: start_hero | all_heroes"),
			widgets.NewDimmedLabelWidget(theme, "Example: StartingGold|all_heroes|5000|"),
			multilineEditor(theme, &this.bonusesEdit, "PresetType|ReceiverFilter|Param|Param2"),
		}),

		widgets.NewSectionWidget(theme, "Banned items", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One item SID per line."),
			multilineEditor(theme, &this.bannedItemsEdit, "item_sid"),
		}),

		widgets.NewSectionWidget(theme, "Banned spells", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One spell SID per line."),
			multilineEditor(theme, &this.bannedMagicsEdit, "spell_sid"),
		}),

		widgets.NewSectionWidget(theme, "Guard value overrides", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One override per line. Format: sid=guardValue"),
			widgets.NewDimmedLabelWidget(theme, "Example: archery_range=350"),
			multilineEditor(theme, &this.valueOverrideEdit, "sid=guardValue"),
		}),
	}
	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

func (this *BonusesPanel) LoadFromState() {
	settings := this.state.GetSettingsFile()
	this.bonusesEdit.SetText(settings.BonusesJson)
	this.bannedItemsEdit.SetText(settings.BannedItems)
	this.bannedMagicsEdit.SetText(settings.BannedMagics)
	this.valueOverrideEdit.SetText(settings.ValueOverridesText)
}

func (this *BonusesPanel) SaveToState() {
	this.state.UpdateState(func(settings *models.SettingsFile) {
		settings.BonusesJson = normaliseLines(this.bonusesEdit.Text())
		settings.BannedItems = normaliseLines(this.bannedItemsEdit.Text())
		settings.BannedMagics = normaliseLines(this.bannedMagicsEdit.Text())
		settings.ValueOverridesText = normaliseLines(this.valueOverrideEdit.Text())
	})
}

// multilineEditor wraps NewTextboxWidget with a sensible minimum height
// so multi-line content stays visible inside the scrolling section.
func multilineEditor(theme *material.Theme, editor *widget.Editor, hint string) layout.Widget {
	inner := widgets.NewTextboxWidget(theme, editor, hint)
	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(96))
		return inner(gtx)
	}
}

// normaliseLines trims trailing whitespace on every line, drops trailing
// blank lines, and converts \r\n to \n so the persisted form is stable.
func normaliseLines(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	// Drop trailing empty lines.
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}
