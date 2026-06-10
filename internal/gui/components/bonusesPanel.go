package components

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
)

type BonusesPanel struct {
	bonusesEdit       widget.Editor
	bannedItemsEdit   widget.Editor
	bannedMagicsEdit  widget.Editor
	valueOverrideEdit widget.Editor

	addBonusBtn      widget.Clickable
	pickItemsBtn     widget.Clickable
	pickSpellsBtn    widget.Clickable
	pickOverridesBtn widget.Clickable

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
			widgets.NewGoldButtonWidget(theme, "Add bonus…", &this.addBonusBtn, false),
			multilineEditor(theme, &this.bonusesEdit, "PresetType|ReceiverFilter|Param|Param2"),
		}),

		widgets.NewSectionWidget(theme, "Banned items", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One item SID per line."),
			widgets.NewButtonWidget(theme, "Pick items…", &this.pickItemsBtn, false),
			multilineEditor(theme, &this.bannedItemsEdit, "item_sid"),
		}),

		widgets.NewSectionWidget(theme, "Banned spells", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One spell SID per line."),
			widgets.NewButtonWidget(theme, "Pick spells…", &this.pickSpellsBtn, false),
			multilineEditor(theme, &this.bannedMagicsEdit, "spell_sid"),
		}),

		widgets.NewSectionWidget(theme, "Guard value overrides", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One override per line. Format: sid=guardValue"),
			widgets.NewDimmedLabelWidget(theme, "Example: archery_range=350"),
			widgets.NewButtonWidget(theme, "Pick overrides…", &this.pickOverridesBtn, false),
			multilineEditor(theme, &this.valueOverrideEdit, "sid=guardValue"),
		}),
	}
	return func(gtx layout.Context) layout.Dimensions {
		this.processPickerClicks(gtx)
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

// processPickerClicks opens the relevant picker dialog when one of the section
// buttons is pressed and appends the chosen values back into the editors.
func (this *BonusesPanel) processPickerClicks(gtx layout.Context) {
	opener := this.state.Dialogs().Open

	if this.addBonusBtn.Clicked(gtx) {
		existing := config_inner.ParseBonusesJSON(this.bonusesEdit.Text())
		opener(NewBonusPickerDialog(existing, opener, func(entries []config_inner.BonusEntry) {
			lines := make([]string, 0, len(entries))
			for _, entry := range entries {
				lines = append(lines, entry.String())
			}
			this.appendLines(&this.bonusesEdit, lines)
		}))
	}
	if this.pickItemsBtn.Clicked(gtx) {
		excluded := splitNonEmptyLines(this.bannedItemsEdit.Text())
		opener(NewItemPickerDialog("Ban Items", excluded, func(ids []string) {
			this.appendLines(&this.bannedItemsEdit, ids)
		}))
	}
	if this.pickSpellsBtn.Clicked(gtx) {
		excluded := splitNonEmptyLines(this.bannedMagicsEdit.Text())
		opener(NewSpellPickerDialog(excluded, false, func(ids []string, _ bool) {
			this.appendLines(&this.bannedMagicsEdit, ids)
		}))
	}
	if this.pickOverridesBtn.Clicked(gtx) {
		excluded := overrideSids(this.valueOverrideEdit.Text())
		opener(NewValueOverridePickerDialog(excluded, func(lines []string) {
			this.appendLines(&this.valueOverrideEdit, lines)
		}))
	}
}

// appendLines appends new, non-duplicate lines to the editor and persists the
// result to state.
func (this *BonusesPanel) appendLines(editor *widget.Editor, newLines []string) {
	lines := splitNonEmptyLines(editor.Text())
	seen := make(map[string]bool, len(lines))
	for _, line := range lines {
		seen[line] = true
	}
	for _, line := range newLines {
		if line == "" || seen[line] {
			continue
		}
		seen[line] = true
		lines = append(lines, line)
	}
	editor.SetText(strings.Join(lines, "\n"))
	this.SaveToState()
}

func (this *BonusesPanel) LoadFromState() {
	settings := this.state.GetStateData()
	this.bonusesEdit.SetText(settings.BonusesJSON)
	this.bannedItemsEdit.SetText(settings.BannedItems)
	this.bannedMagicsEdit.SetText(settings.BannedMagics)
	this.valueOverrideEdit.SetText(settings.ValueOverridesText)
}

func (this *BonusesPanel) SaveToState() {
	this.state.UpdateState(func(settings *models.EditorStateModel) {
		settings.BonusesJSON = normaliseLines(this.bonusesEdit.Text())
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

// splitNonEmptyLines returns the trimmed, non-empty lines of text.
func splitNonEmptyLines(text string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// overrideSids extracts the SID portion (before '=') of each existing override
// line so the value-override picker can hide already-overridden SIDs.
func overrideSids(text string) []string {
	lines := splitNonEmptyLines(text)
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if idx := strings.IndexByte(line, '='); idx >= 0 {
			out = append(out, strings.TrimSpace(line[:idx]))
		} else {
			out = append(out, line)
		}
	}
	return out
}
