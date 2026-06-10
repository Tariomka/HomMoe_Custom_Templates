package components

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
)

// Category dot colours, matching the C# editor's brushes.
var (
	dotMovement = color.NRGBA{R: 0x64, G: 0x95, B: 0xED, A: 0xFF} // cornflower blue
	dotCombat   = color.NRGBA{R: 0xCD, G: 0x5C, B: 0x5C, A: 0xFF} // indian red
	dotMagic    = color.NRGBA{R: 0x93, G: 0x70, B: 0xDB, A: 0xFF} // medium purple
	dotSet      = color.NRGBA{R: 0xBA, G: 0x55, B: 0xD3, A: 0xFF} // medium orchid
	dotResource = color.NRGBA{R: 0xDA, G: 0xA5, B: 0x20, A: 0xFF} // goldenrod
	dotDefault  = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF} // gray
)

// BonusesPanel mirrors the parallel C# editor's bonuses & bans tab: bonuses
// and bans are shown as read-only lists with human-readable names; entries
// are added through picker dialogs and removed with per-row buttons. Only
// the guard-value overrides stay free-text (matching the C# editor).
type BonusesPanel struct {
	bonuses      []config_inner.BonusEntry
	bannedItems  []string
	bannedMagics []string

	bonusRemoveBtns []widget.Clickable
	itemRemoveBtns  []widget.Clickable
	magicRemoveBtns []widget.Clickable

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
	return func(gtx layout.Context) layout.Dimensions {
		this.processClicks(gtx)
		widgetsList := this.buildWidgets(theme)
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

// buildWidgets assembles the section list for the current frame; rebuilt every
// frame so the rows stay in sync with the underlying entry slices.
func (this *BonusesPanel) buildWidgets(theme *material.Theme) []layout.Widget {
	bonusRows := []layout.Widget{
		widgets.NewGoldButtonWidget(theme, "＋ Add bonus…", &this.addBonusBtn, false),
	}
	if len(this.bonuses) == 0 {
		bonusRows = append(bonusRows, widgets.NewDimmedLabelWidget(theme, "(no bonuses)"))
	}
	for i, entry := range this.bonuses {
		bonusRows = append(bonusRows, this.entryRow(theme,
			bonusDotColor(entry.PresetType),
			bonusDisplayName(entry),
			bonusReceiverLabel(entry),
			&this.bonusRemoveBtns[i],
		))
	}

	itemRows := []layout.Widget{
		widgets.NewGoldButtonWidget(theme, "＋ Add banned item…", &this.pickItemsBtn, false),
	}
	if len(this.bannedItems) == 0 {
		itemRows = append(itemRows, widgets.NewDimmedLabelWidget(theme, "(no banned items)"))
	}
	for i, sid := range this.bannedItems {
		name, category := bannedItemLabel(sid)
		itemRows = append(itemRows, this.entryRow(theme, banCategoryColor(category), name, category, &this.itemRemoveBtns[i]))
	}

	spellRows := []layout.Widget{
		widgets.NewGoldButtonWidget(theme, "＋ Add banned spell…", &this.pickSpellsBtn, false),
	}
	if len(this.bannedMagics) == 0 {
		spellRows = append(spellRows, widgets.NewDimmedLabelWidget(theme, "(no banned spells)"))
	}
	for i, sid := range this.bannedMagics {
		name, school := bannedSpellLabel(sid)
		spellRows = append(spellRows, this.entryRow(theme, dotMagic, name, school, &this.magicRemoveBtns[i]))
	}

	return []layout.Widget{
		widgets.NewWarningBannerWidget(theme, "EXPERIMENTAL — Effects only apply on generation."),
		widgets.NewSectionWidget(theme, "Game start bonuses", bonusRows),
		widgets.NewSectionWidget(theme, "Banned items", itemRows),
		widgets.NewSectionWidget(theme, "Banned spells", spellRows),
		widgets.NewSectionWidget(theme, "Guard value overrides", []layout.Widget{
			widgets.NewDimmedLabelWidget(theme, "One override per line. Format: sid=guardValue"),
			widgets.NewDimmedLabelWidget(theme, "Example: archery_range=350"),
			widgets.NewButtonWidget(theme, "＋ Pick overrides…", &this.pickOverridesBtn, false),
			multilineEditor(theme, &this.valueOverrideEdit, "sid=guardValue"),
		}),
	}
}

// entryRow renders one read-only list row: category dot, primary label, dim
// trailing label, and a remove button.
func (this *BonusesPanel) entryRow(theme *material.Theme, dot color.NRGBA, name, trailing string, removeBtn *widget.Clickable) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(dotWidget(dot)),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(theme, name)
					label.Color = themes.ColorText
					label.TextSize = unit.Sp(13)
					return label.Layout(gtx)
				}),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if trailing == "" {
						return layout.Dimensions{}
					}
					label := material.Body2(theme, trailing)
					label.Color = themes.ColorTextDim
					label.TextSize = unit.Sp(11)
					return label.Layout(gtx)
				}),
				layout.Rigid(widgets.NewButtonWidget(theme, "✕", removeBtn, false)),
			)
		})
	}
}

// dotWidget draws the small filled category circle.
func dotWidget(col color.NRGBA) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := gtx.Dp(unit.Dp(8))
		paint.FillShape(gtx.Ops, col, clip.Ellipse{Max: image.Pt(size, size)}.Op(gtx.Ops))
		return layout.Dimensions{Size: image.Pt(size, size)}
	}
}

// processClicks handles add-dialog launches and per-row removals.
func (this *BonusesPanel) processClicks(gtx layout.Context) {
	opener := this.state.Dialogs().Open

	if this.addBonusBtn.Clicked(gtx) {
		opener(NewBonusPickerDialog(this.bonuses, opener, func(entries []config_inner.BonusEntry) {
			this.bonuses = append(this.bonuses, entries...)
			this.syncRemoveButtons()
			this.SaveToState()
		}))
	}
	if this.pickItemsBtn.Clicked(gtx) {
		opener(NewItemPickerDialog("Ban Items", this.bannedItems, func(ids []string) {
			this.bannedItems = appendUnique(this.bannedItems, ids)
			this.syncRemoveButtons()
			this.SaveToState()
		}))
	}
	if this.pickSpellsBtn.Clicked(gtx) {
		opener(NewSpellPickerDialog(this.bannedMagics, false, func(ids []string, _ bool) {
			this.bannedMagics = appendUnique(this.bannedMagics, ids)
			this.syncRemoveButtons()
			this.SaveToState()
		}))
	}
	if this.pickOverridesBtn.Clicked(gtx) {
		excluded := overrideSids(this.valueOverrideEdit.Text())
		opener(NewValueOverridePickerDialog(excluded, func(lines []string) {
			this.appendOverrideLines(lines)
		}))
	}

	if index := clickedIndex(gtx, this.bonusRemoveBtns); index >= 0 {
		this.bonuses = append(this.bonuses[:index:index], this.bonuses[index+1:]...)
		this.syncRemoveButtons()
		this.SaveToState()
	}
	if index := clickedIndex(gtx, this.itemRemoveBtns); index >= 0 {
		this.bannedItems = append(this.bannedItems[:index:index], this.bannedItems[index+1:]...)
		this.syncRemoveButtons()
		this.SaveToState()
	}
	if index := clickedIndex(gtx, this.magicRemoveBtns); index >= 0 {
		this.bannedMagics = append(this.bannedMagics[:index:index], this.bannedMagics[index+1:]...)
		this.syncRemoveButtons()
		this.SaveToState()
	}
}

// syncRemoveButtons resizes the per-row clickable pools to the list lengths.
func (this *BonusesPanel) syncRemoveButtons() {
	this.bonusRemoveBtns = resizeClickables(this.bonusRemoveBtns, len(this.bonuses))
	this.itemRemoveBtns = resizeClickables(this.itemRemoveBtns, len(this.bannedItems))
	this.magicRemoveBtns = resizeClickables(this.magicRemoveBtns, len(this.bannedMagics))
}

// appendOverrideLines appends new, non-duplicate lines to the override editor
// and persists the result to state.
func (this *BonusesPanel) appendOverrideLines(newLines []string) {
	lines := splitNonEmptyLines(this.valueOverrideEdit.Text())
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
	this.valueOverrideEdit.SetText(strings.Join(lines, "\n"))
	this.SaveToState()
}

func (this *BonusesPanel) LoadFromState() {
	settings := this.state.GetStateData()
	this.bonuses = config_inner.ParseBonusesJSON(settings.BonusesJSON)
	this.bannedItems = splitNonEmptyLines(settings.BannedItems)
	this.bannedMagics = splitNonEmptyLines(settings.BannedMagics)
	this.valueOverrideEdit.SetText(settings.ValueOverridesText)
	this.syncRemoveButtons()
}

func (this *BonusesPanel) SaveToState() {
	this.state.UpdateState(func(settings *models.EditorStateModel) {
		settings.BonusesJSON = config_inner.SerializeBonuses(this.bonuses)
		settings.BannedItems = strings.Join(this.bannedItems, "\n")
		settings.BannedMagics = strings.Join(this.bannedMagics, "\n")
		settings.ValueOverridesText = normaliseLines(this.valueOverrideEdit.Text())
	})
}

// ── display helpers ─────────────────────────────────────────────────────────

// bonusDisplayName composes the human-readable label for a bonus entry,
// mirroring the C# BonusEntry.DisplayName.
func bonusDisplayName(entry config_inner.BonusEntry) string {
	switch entry.PresetType {
	case config_inner.BonusTownPortalFree:
		return "Town Portal (free)"
	case config_inner.BonusSpell:
		if entry.Param2 == "1" {
			return "Spell (free): " + spellLabel(entry.Param)
		}
		return "Spell: " + spellLabel(entry.Param)
	case config_inner.BonusUnitMultiplier:
		return "Unit multiplier ×" + entry.Param
	case config_inner.BonusMovementBonus:
		return "Movement bonus +" + entry.Param
	case config_inner.BonusStartingItem:
		return "Starting item: " + itemLabel(entry.Param)
	case config_inner.BonusStartingGold:
		return "Starting gold: " + entry.Param
	case config_inner.BonusStartingGems:
		return "Starting gems: " + entry.Param
	case config_inner.BonusStartingCrystals:
		return "Starting crystals: " + entry.Param
	case config_inner.BonusStartingMercury:
		return "Starting mercury: " + entry.Param
	case config_inner.BonusStartingWood:
		return "Starting wood: " + entry.Param
	case config_inner.BonusStartingOre:
		return "Starting ore: " + entry.Param
	}
	return entry.String()
}

// bonusReceiverLabel is the dim trailing text; hidden for resource bonuses,
// matching the C# ShowReceiverLabel behaviour.
func bonusReceiverLabel(entry config_inner.BonusEntry) string {
	if isResourceBonus(entry.PresetType) {
		return ""
	}
	if entry.ReceiverFilter == "all_heroes" {
		return "all heroes"
	}
	return "start hero"
}

// bonusDotColor matches the C# BonusEntry.DotBrush colour coding.
func bonusDotColor(typ config_inner.BonusPresetType) color.NRGBA {
	switch typ {
	case config_inner.BonusTownPortalFree, config_inner.BonusSpell:
		return dotMagic
	case config_inner.BonusUnitMultiplier:
		return dotCombat
	case config_inner.BonusMovementBonus:
		return dotMovement
	case config_inner.BonusStartingItem:
		return dotSet
	default:
		return dotResource
	}
}

// banCategoryColor matches the C# BanEntry.CategoryBrush colour coding.
func banCategoryColor(category string) color.NRGBA {
	switch category {
	case "Movement":
		return dotMovement
	case "Diplomacy":
		return dotResource
	case "Combat":
		return dotCombat
	case "Magic", "Spell":
		return dotMagic
	case "Set":
		return dotSet
	}
	return dotDefault
}

// spellLabel resolves a spell SID to its display name, with a generic
// sentence-case fallback for unknown SIDs.
func spellLabel(sid string) string {
	if spell, ok := constants.FindSpell(sid); ok {
		return spell.Name
	}
	return constants.SidToDisplayName(sid)
}

// itemLabel resolves an artifact SID to its display name, with a generic
// sentence-case fallback for unknown SIDs.
func itemLabel(sid string) string {
	if item, ok := constants.FindBannableItem(sid); ok {
		return item.Name
	}
	return constants.SidToDisplayName(sid)
}

// bannedItemLabel returns the display name and category for a banned artifact.
func bannedItemLabel(sid string) (name, category string) {
	if item, ok := constants.FindBannableItem(sid); ok {
		return item.Name, item.Category
	}
	return constants.SidToDisplayName(sid), "Misc"
}

// bannedSpellLabel returns the display name and school label for a banned spell.
func bannedSpellLabel(sid string) (name, school string) {
	if spell, ok := constants.FindSpell(sid); ok {
		label := constants.SpellSchoolDisplayNames[spell.School]
		if label == "" {
			label = spell.School
		}
		return spell.Name, label
	}
	return constants.SidToDisplayName(sid), "Spell"
}

// ── small helpers ───────────────────────────────────────────────────────────

// clickedIndex returns the index of the first clicked button, or -1.
func clickedIndex(gtx layout.Context, buttons []widget.Clickable) int {
	for i := range buttons {
		if buttons[i].Clicked(gtx) {
			return i
		}
	}
	return -1
}

// resizeClickables grows or shrinks a clickable pool to the wanted length.
func resizeClickables(buttons []widget.Clickable, length int) []widget.Clickable {
	for len(buttons) < length {
		buttons = append(buttons, widget.Clickable{})
	}
	return buttons[:length]
}

// appendUnique appends the ids not already present in values.
func appendUnique(values, ids []string) []string {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		seen[value] = true
	}
	for _, id := range ids {
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		values = append(values, id)
	}
	return values
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
