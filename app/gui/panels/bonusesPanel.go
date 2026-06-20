package panels

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
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
)

// Category dot colors come from the central theme palette.
var (
	dotMovement = themes.ColorDotMovement
	dotCombat   = themes.ColorDotCombat
	dotMagic    = themes.ColorDotMagic
	dotSet      = themes.ColorDotSet
	dotResource = themes.ColorDotResource
	dotDefault  = themes.ColorDotDefault
)

// BonusesPanel mirrors the parallel C# editor's bonuses & bans tab: bonuses,
// bans and guard value overrides are shown as read-only lists with
// human-readable names; entries are added through picker dialogs and removed
// with per-row buttons.
type BonusesPanel struct {
	bonuses        []config_inner.BonusEntry
	bannedItems    []string
	bannedMagics   []string
	valueOverrides []string

	bonusRemoveBtns    []widget.Clickable
	itemRemoveBtns     []widget.Clickable
	magicRemoveBtns    []widget.Clickable
	overrideRemoveBtns []widget.Clickable

	addBonusBtn      widget.Clickable
	pickItemsBtn     widget.Clickable
	pickSpellsBtn    widget.Clickable
	pickOverridesBtn widget.Clickable

	scroll widget.List

	state *drivers.State
}

func NewBonusesPanel(state *drivers.State) *BonusesPanel {
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
		widgets.NewGoldButtonWidget(theme, "+ Add bonus…", &this.addBonusBtn, false),
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
		widgets.NewGoldButtonWidget(theme, "+ Add banned item…", &this.pickItemsBtn, false),
	}
	if len(this.bannedItems) == 0 {
		itemRows = append(itemRows, widgets.NewDimmedLabelWidget(theme, "(no banned items)"))
	}
	for i, sid := range this.bannedItems {
		name, category := bannedItemLabel(sid)
		itemRows = append(itemRows, this.entryRow(theme, banCategoryColor(category), name, category, &this.itemRemoveBtns[i]))
	}

	spellRows := []layout.Widget{
		widgets.NewGoldButtonWidget(theme, "+ Add banned spell…", &this.pickSpellsBtn, false),
	}
	if len(this.bannedMagics) == 0 {
		spellRows = append(spellRows, widgets.NewDimmedLabelWidget(theme, "(no banned spells)"))
	}
	for i, sid := range this.bannedMagics {
		name, school := bannedSpellLabel(sid)
		spellRows = append(spellRows, this.entryRow(theme, dotMagic, name, school, &this.magicRemoveBtns[i]))
	}

	overrideRows := []layout.Widget{
		widgets.NewGoldButtonWidget(theme, "+ Add override…", &this.pickOverridesBtn, false),
	}
	if len(this.valueOverrides) == 0 {
		overrideRows = append(overrideRows, widgets.NewDimmedLabelWidget(theme, "(no overrides)"))
	}
	for i, line := range this.valueOverrides {
		name, value := overrideLabel(line)
		overrideRows = append(overrideRows, this.entryRow(theme, dotResource, name, value, &this.overrideRemoveBtns[i]))
	}

	return []layout.Widget{
		// widgets.NewWarningBannerWidget(theme, "EXPERIMENTAL — Effects only apply on generation."),
		widgets.NewSectionWidget(theme, "Game start bonuses", bonusRows),
		widgets.NewSectionWidget(theme, "Banned items", itemRows),
		widgets.NewSectionWidget(theme, "Banned spells", spellRows),
		widgets.NewSectionWidget(theme, "Guard value overrides", overrideRows),
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
		opener(dialogs.NewBonusPickerDialog(this.bonuses, opener, func(entries []config_inner.BonusEntry) {
			this.bonuses = append(this.bonuses, entries...)
			this.syncRemoveButtons()
			this.SaveToState()
		}))
	}
	if this.pickItemsBtn.Clicked(gtx) {
		opener(dialogs.NewItemPickerDialog("Ban Items", this.bannedItems, func(ids []string) {
			this.bannedItems = appendUnique(this.bannedItems, ids)
			this.syncRemoveButtons()
			this.SaveToState()
		}))
	}
	if this.pickSpellsBtn.Clicked(gtx) {
		opener(dialogs.NewSpellPickerDialog(this.bannedMagics, false, func(ids []string, _ bool) {
			this.bannedMagics = appendUnique(this.bannedMagics, ids)
			this.syncRemoveButtons()
			this.SaveToState()
		}))
	}
	if this.pickOverridesBtn.Clicked(gtx) {
		excluded := overrideSids(this.valueOverrides)
		opener(dialogs.NewValueOverridePickerDialog(excluded, func(lines []string) {
			this.valueOverrides = appendUnique(this.valueOverrides, lines)
			this.syncRemoveButtons()
			this.SaveToState()
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
	if index := clickedIndex(gtx, this.overrideRemoveBtns); index >= 0 {
		this.valueOverrides = append(this.valueOverrides[:index:index], this.valueOverrides[index+1:]...)
		this.syncRemoveButtons()
		this.SaveToState()
	}
}

// syncRemoveButtons resizes the per-row clickable pools to the list lengths.
func (this *BonusesPanel) syncRemoveButtons() {
	this.bonusRemoveBtns = resizeClickables(this.bonusRemoveBtns, len(this.bonuses))
	this.itemRemoveBtns = resizeClickables(this.itemRemoveBtns, len(this.bannedItems))
	this.magicRemoveBtns = resizeClickables(this.magicRemoveBtns, len(this.bannedMagics))
	this.overrideRemoveBtns = resizeClickables(this.overrideRemoveBtns, len(this.valueOverrides))
}

func (this *BonusesPanel) LoadFromState() {
	settings := this.state.GetStateData()
	this.bonuses = config_inner.DeserializeBonusEntries(settings.BonusesJSON)
	this.bannedItems = splitNonEmptyLines(settings.BannedItems)
	this.bannedMagics = splitNonEmptyLines(settings.BannedMagics)
	this.valueOverrides = splitNonEmptyLines(settings.ValueOverridesText)
	this.syncRemoveButtons()
}

func (this *BonusesPanel) SaveToState() {
	this.state.UpdateState(func(settings *dtos.EditorStateDto) {
		settings.BonusesJSON = config_inner.SerializeBonusEntries(this.bonuses)
		settings.BannedItems = strings.Join(this.bannedItems, "\n")
		settings.BannedMagics = strings.Join(this.bannedMagics, "\n")
		settings.ValueOverridesText = strings.Join(this.valueOverrides, "\n")
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
	if entry.PresetType.IsResource() {
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

// overrideSids extracts the SID portion (before '=') of each override line so
// the value-override picker can hide already-overridden SIDs.
func overrideSids(lines []string) []string {
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

// overrideLabel splits a "sid=guardValue" line into a display name and a dim
// trailing label describing the override value.
func overrideLabel(line string) (name, trailing string) {
	sid := line
	value := ""
	if idx := strings.IndexByte(line, '='); idx >= 0 {
		sid = strings.TrimSpace(line[:idx])
		value = strings.TrimSpace(line[idx+1:])
	}
	if value == "" {
		return constants.SidToDisplayName(sid), sid
	}
	return constants.SidToDisplayName(sid), "guard value: " + value
}
