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
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
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

type BonusesPanel struct {
	bonuses        []config.BonusEntry
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
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList),
			func(gtx layout.Context, index int) layout.Dimensions { return widgetsList[index](gtx) })
	}
}

// buildWidgets assembles the section list for the current frame; rebuilt every
// frame so the rows stay in sync with the underlying entry slices.
func (this *BonusesPanel) buildWidgets(theme *material.Theme) []layout.Widget {
	return []layout.Widget{
		// widgets.NewWarningBannerWidget(theme, "EXPERIMENTAL - Effects only apply on generation."),
		widgets.NewHorizontallySplitWidget(theme,
			func(theme *material.Theme) layout.Widget {
				return func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(widgets.NewSectionWidget(theme, "Game start bonuses",
							this.getStartBonusesWidgets(theme))),
						layout.Rigid(widgets.NewSectionWidget(theme, "Guard value overrides",
							this.getValueOverridesWidgets(theme))))
				}
			},
			func(theme *material.Theme) layout.Widget {
				return func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(widgets.NewSectionWidget(theme, "Banned items",
							this.getBannedItemsWidgets(theme))),
						layout.Rigid(widgets.NewSectionWidget(theme, "Banned spells",
							this.getBannedSpellsWidgets(theme))))
				}
			}),
	}
}

func (this *BonusesPanel) getStartBonusesWidgets(theme *material.Theme) []layout.Widget {
	bonusRows := []layout.Widget{
		widgets.NewBrightButtonLargeWidget(theme, "+ Add bonus...", &this.addBonusBtn, false),
	}
	if len(this.bonuses) == 0 {
		return append(bonusRows, widgets.NewDimmedLabelWidget(theme, "(no bonuses)"))
	}

	for i, entry := range this.bonuses {
		bonusRows = append(bonusRows, this.getEntryRowWidget(theme,
			bonusDotColor(entry.PresetType),
			bonusDisplayName(entry),
			bonusReceiverLabel(entry),
			&this.bonusRemoveBtns[i]))
	}
	return bonusRows
}

func (this *BonusesPanel) getBannedItemsWidgets(theme *material.Theme) []layout.Widget {
	itemRows := []layout.Widget{
		widgets.NewBrightButtonLargeWidget(theme, "+ Add banned item...", &this.pickItemsBtn, false),
	}
	if len(this.bannedItems) == 0 {
		return append(itemRows, widgets.NewDimmedLabelWidget(theme, "(no banned items)"))
	}

	for i, sid := range this.bannedItems {
		name, category := bannedItemLabel(sid)
		itemRows = append(itemRows,
			this.getEntryRowWidget(theme, banCategoryColor(category), name, category, &this.itemRemoveBtns[i]))
	}
	return itemRows
}

func (this *BonusesPanel) getBannedSpellsWidgets(theme *material.Theme) []layout.Widget {
	spellRows := []layout.Widget{
		widgets.NewBrightButtonLargeWidget(theme, "+ Add banned spell...", &this.pickSpellsBtn, false),
	}
	if len(this.bannedMagics) == 0 {
		return append(spellRows, widgets.NewDimmedLabelWidget(theme, "(no banned spells)"))
	}

	for i, sid := range this.bannedMagics {
		name, school := bannedSpellLabel(sid)
		spellRows = append(spellRows, this.getEntryRowWidget(theme,
			constants.GetSpellSchoolColorFromDisplayName(school), name, school, &this.magicRemoveBtns[i]),
		)
	}
	return spellRows
}

func (this *BonusesPanel) getValueOverridesWidgets(theme *material.Theme) []layout.Widget {
	overrideRows := []layout.Widget{
		widgets.NewBrightButtonLargeWidget(theme, "+ Add override...", &this.pickOverridesBtn, false),
	}
	if len(this.valueOverrides) == 0 {
		return append(overrideRows, widgets.NewDimmedLabelWidget(theme, "(no overrides)"))
	}

	for i, line := range this.valueOverrides {
		name, value := overrideLabel(line)
		overrideRows = append(overrideRows,
			this.getEntryRowWidget(theme, dotResource, name, value, &this.overrideRemoveBtns[i]))
	}
	return overrideRows
}

// getEntryRowWidget renders one read-only list row: category dot, primary label, dim
// trailing label, and a remove button.
func (this *BonusesPanel) getEntryRowWidget(
	theme *material.Theme,
	dotColor color.NRGBA,
	name, trailing string,
	removeBtn *widget.Clickable) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(constants.DefaultRoundnessLarge)
						paint.FillShape(gtx.Ops, dotColor, clip.Ellipse{Max: image.Pt(size, size)}.Op(gtx.Ops))
						return layout.Dimensions{Size: image.Pt(size, size)}
					}),
					widgets.NewDefaultComponentSpacer(),
					layout.Rigid(widgets.NewLabelBigWidget(theme, name, themes.ColorText)),
					widgets.NewDefaultComponentSpacer(),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if trailing == "" {
							return layout.Dimensions{}
						}

						return widgets.NewDimmedLabelWidget(theme, trailing)(gtx)
					}),
					layout.Rigid(widgets.NewButtonWidget(theme, "✕", removeBtn, false)),
				)
			},
		)
	}
}

// processClicks handles add-dialog launches and per-row removals.
func (this *BonusesPanel) processClicks(gtx layout.Context) {
	opener := this.state.GetDialogHost().Open

	if this.addBonusBtn.Clicked(gtx) {
		opener(dialogs.NewBonusPickerDialog(this.bonuses, opener, func(entries []config.BonusEntry) {
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
	this.bonuses = settings.BonusesJSON
	this.bannedItems = splitNonEmptyLines(settings.BannedItems)
	this.bannedMagics = splitNonEmptyLines(settings.BannedMagics)
	this.valueOverrides = splitNonEmptyLines(settings.ValueOverridesText)
	this.syncRemoveButtons()
}

func (this *BonusesPanel) SaveToState() {
	this.state.UpdateState(func(settings *dtos.EditorStateDto) {
		settings.BonusesJSON = this.bonuses
		settings.BannedItems = strings.Join(this.bannedItems, "\n")
		settings.BannedMagics = strings.Join(this.bannedMagics, "\n")
		settings.ValueOverridesText = strings.Join(this.valueOverrides, "\n")
	})
}

// ── display helpers ─────────────────────────────────────────────────────────

// bonusDisplayName composes the human-readable label for a bonus entry.
func bonusDisplayName(entry config.BonusEntry) string {
	switch entry.PresetType {
	case config.BonusTownPortalFree:
		return "Town Portal (free)"
	case config.BonusSpell:
		if entry.Param2 == "1" {
			return "Spell (free): " + spellLabel(entry.Param)
		}
		return "Spell: " + spellLabel(entry.Param)
	case config.BonusUnitMultiplier:
		return "Unit multiplier x" + entry.Param
	case config.BonusMovementBonus:
		return "Movement bonus +" + entry.Param
	case config.BonusStartingItem:
		return "Starting item: " + itemLabel(entry.Param)
	case config.BonusStartingGold:
		return "Starting gold: " + entry.Param
	case config.BonusStartingGems:
		return "Starting gems: " + entry.Param
	case config.BonusStartingCrystals:
		return "Starting crystals: " + entry.Param
	case config.BonusStartingMercury:
		return "Starting mercury: " + entry.Param
	case config.BonusStartingWood:
		return "Starting wood: " + entry.Param
	case config.BonusStartingOre:
		return "Starting ore: " + entry.Param
	}
	return ""
}

// bonusReceiverLabel is the dim trailing text; hidden for resource bonuses.
func bonusReceiverLabel(entry config.BonusEntry) string {
	if entry.PresetType.IsResource() {
		return ""
	}

	if entry.ReceiverFilter == "all_heroes" {
		return "all heroes"
	}

	return "start hero"
}

// bonusDotColor matches the C# BonusEntry.DotBrush color coding.
func bonusDotColor(typ config.BonusPresetType) color.NRGBA {
	switch typ {
	case config.BonusTownPortalFree, config.BonusSpell:
		return dotMagic
	case config.BonusUnitMultiplier:
		return dotCombat
	case config.BonusMovementBonus:
		return dotMovement
	case config.BonusStartingItem:
		return dotSet
	default:
		return dotResource
	}
}

// banCategoryColor matches the C# BanEntry.CategoryBrush color coding.
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
	for line := range strings.SplitSeq(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
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
		if before, _, ok := strings.Cut(line, "="); ok {
			out = append(out, strings.TrimSpace(before))
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
	if before, after, ok := strings.Cut(line, "="); ok {
		sid = strings.TrimSpace(before)
		value = strings.TrimSpace(after)
	}
	if value == "" {
		return constants.SidToDisplayName(sid), sid
	}

	return constants.SidToDisplayName(sid), "guard value: " + value
}
