package dialogs

import (
	"image"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
)

// bonusTypeOption pairs a friendly dropdown label with its preset type.
type bonusTypeOption struct {
	label string
	typ   config_inner.BonusPresetType
}

var bonusTypeOptions = []bonusTypeOption{
	{"Free Town Portal", config_inner.BonusTownPortalFree},
	{"Spell", config_inner.BonusSpell},
	{"Unit Multiplier", config_inner.BonusUnitMultiplier},
	{"Movement Bonus", config_inner.BonusMovementBonus},
	{"Starting Item", config_inner.BonusStartingItem},
	{"Starting Gold", config_inner.BonusStartingGold},
	{"Starting Gems", config_inner.BonusStartingGems},
	{"Starting Crystals", config_inner.BonusStartingCrystals},
	{"Starting Mercury", config_inner.BonusStartingMercury},
	{"Starting Wood", config_inner.BonusStartingWood},
	{"Starting Ore", config_inner.BonusStartingOre},
}

var bonusReceiverOptions = []string{"start_hero", "all_heroes"}

var bonusResourceDefaults = map[config_inner.BonusPresetType]string{
	config_inner.BonusStartingGold:     "10000",
	config_inner.BonusStartingGems:     "15",
	config_inner.BonusStartingCrystals: "15",
	config_inner.BonusStartingMercury:  "15",
	config_inner.BonusStartingWood:     "20",
	config_inner.BonusStartingOre:      "20",
}

// BonusPickerDialog composes one or more game-start bonuses. For the Spell
// type it keeps a removable list of picked spells (the spell picker appends to
// it) and emits one bonus per spell; for StartingItem it can launch the item
// picker (via opener). Implements widgets.Dialog.
type BonusPickerDialog struct {
	existingKeys     map[string]bool
	existingSpellIds []string
	opener           interfaces.DialogOpener
	onApply          func(entries []config_inner.BonusEntry)

	typeDropdown     *components.DropdownSelector
	receiverDropdown *components.DropdownSelector

	selectedSpells  []string
	spellRemoveBtns []widget.Clickable
	spellScroll     widget.List
	makeFree        widget.Bool
	pickSpellBtn    widget.Clickable

	multiplierEdit widget.Editor
	movementEdit   widget.Editor

	itemEdit    widget.Editor
	pickItemBtn widget.Clickable

	resourceEdit widget.Editor

	addBtn    widget.Clickable
	cancelBtn widget.Clickable

	errorText  string
	pending    []config_inner.BonusEntry
	hasPending bool
}

// NewBonusPickerDialog builds the bonus composer. existing is used for
// duplicate detection and to pre-exclude already-chosen spells in the spell
// sub-picker. onApply receives the composed bonus entries.
func NewBonusPickerDialog(existing []config_inner.BonusEntry, opener interfaces.DialogOpener, onApply func(entries []config_inner.BonusEntry)) *BonusPickerDialog {
	keys := make(map[string]bool, len(existing))
	var spellIds []string
	for _, entry := range existing {
		keys[entry.GetHash()] = true
		if entry.PresetType == config_inner.BonusSpell && entry.Param != "" {
			spellIds = append(spellIds, entry.Param)
		}
	}

	labels := make([]string, len(bonusTypeOptions))
	for i, option := range bonusTypeOptions {
		labels[i] = option.label
	}

	dialog := &BonusPickerDialog{
		existingKeys:     keys,
		existingSpellIds: spellIds,
		opener:           opener,
		onApply:          onApply,
		typeDropdown:     components.NewDropdownSelector(labels),
		receiverDropdown: components.NewDropdownSelector(append([]string{}, bonusReceiverOptions...)),
	}
	dialog.spellScroll.Axis = layout.Vertical
	dialog.multiplierEdit.SingleLine = true
	dialog.multiplierEdit.SetText("2")
	dialog.movementEdit.SingleLine = true
	dialog.movementEdit.SetText("0")
	dialog.itemEdit.SingleLine = true
	dialog.itemEdit.ReadOnly = true
	dialog.resourceEdit.SingleLine = true
	return dialog
}

func (this *BonusPickerDialog) Title() string { return "Add Bonus" }

func (this *BonusPickerDialog) PreferredSize() (unit.Dp, unit.Dp) {
	return unit.Dp(460), unit.Dp(420)
}

func (this *BonusPickerDialog) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	if this.handleClicks(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}
	this.handleSubPickers(gtx)

	presetType := this.getSelectedType()

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := widgets.NewLabeledRowWidget(theme, "Type", 90, this.typeDropdown.GetWidget(theme))(gtx)
			if this.typeDropdown.WasUpdated {
				this.applyTypeDefaults(presetType)
			}
			return dims
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if presetType.IsResource() {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(theme, "Receiver", 90, this.receiverDropdown.GetWidget(theme))(gtx)
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Rigid(this.getEditorWidget(theme, presetType)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if this.errorText == "" {
				return layout.Dimensions{}
			}

			return layout.Inset{Top: constants.DefaultPadding}.Layout(gtx,
				widgets.NewLabelWidget(theme, this.errorText, themes.ColorError))
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(12)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, widgets.NewMinimumConstraintsWidget()),
				layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelBtn, false)),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
				layout.Rigid(widgets.NewBrightButtonWidget(theme, "Add", &this.addBtn, false)),
			)
		}),
	), false
}

func (this *BonusPickerDialog) getEditorWidget(theme *material.Theme, presetType config_inner.BonusPresetType) layout.Widget {
	switch presetType {
	case config_inner.BonusTownPortalFree:
		return widgets.NewDimmedLabelWidget(theme, "Grants free town portal to the selected hero(es). No extra parameters.")
	case config_inner.BonusSpell:
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, widgets.NewDimmedLabelWidget(theme, spellCountLabel(len(this.selectedSpells)))),
						layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
						layout.Rigid(widgets.NewButtonWidget(
							theme, "Pick spells...", &this.pickSpellBtn, this.opener == nil)),
					)
				}),
				layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
				layout.Rigid(this.getSpellListWidget(theme)),
				layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
				layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.makeFree, "Make spell(s) free")),
			)
		}
	case config_inner.BonusUnitMultiplier:
		return widgets.NewLabeledRowWidget(theme, "Multiplier", 90,
			widgets.NewTextboxWidget(theme, &this.multiplierEdit, "2", false))
	case config_inner.BonusMovementBonus:
		return widgets.NewLabeledRowWidget(theme, "Movement", 90,
			widgets.NewTextboxWidget(theme, &this.movementEdit, "0", false))
	case config_inner.BonusStartingItem:
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.itemEdit, "use Pick item...", false)),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
				layout.Rigid(widgets.NewButtonWidget(theme, "Pick item", &this.pickItemBtn, this.opener == nil)),
			)
		}
	default: // resources
		return widgets.NewLabeledRowWidget(theme, "Amount", 90,
			widgets.NewTextboxWidget(theme, &this.resourceEdit, "0", false))
	}
}

// handleClicks polls the clickable buttons and returns true if the dialog should close.
func (this *BonusPickerDialog) handleClicks(gtx layout.Context) bool {
	// A multi-selection sub-picker queued several bonuses; emit and close.
	if this.hasPending {
		this.hasPending = false
		if this.onApply != nil {
			this.onApply(this.pending)
		}
		return true
	}

	if this.cancelBtn.Clicked(gtx) {
		return true
	}

	// Poll the per-row remove buttons BEFORE they are laid out - Clickable.Layout
	// consumes the click, so a check after layout never fires.
	for i := range this.spellRemoveBtns {
		if this.spellRemoveBtns[i].Clicked(gtx) {
			this.selectedSpells = append(this.selectedSpells[:i:i], this.selectedSpells[i+1:]...)
			this.spellRemoveBtns = this.spellRemoveBtns[:len(this.selectedSpells)]
			break
		}
	}

	if this.addBtn.Clicked(gtx) {
		if entries, ok := this.buildEntries(); ok {
			fresh := make([]config_inner.BonusEntry, 0, len(entries))
			for _, entry := range entries {
				if !this.existingKeys[entry.GetHash()] {
					fresh = append(fresh, entry)
				}
			}
			if len(fresh) == 0 {
				this.errorText = "That bonus already exists."
			} else {
				if this.onApply != nil {
					this.onApply(fresh)
				}
				return true
			}
		}
	}

	return false
}

// handleSubPickers reacts to the "Pick spell" / "Pick item" buttons by pushing
// the relevant picker onto the dialog stack.
func (this *BonusPickerDialog) handleSubPickers(gtx layout.Context) {
	if this.opener == nil {
		return
	}

	if this.pickSpellBtn.Clicked(gtx) {
		excluded := append(append([]string{}, this.existingSpellIds...), this.selectedSpells...)
		this.opener(NewSpellPickerDialog(excluded, false, func(ids []string, _ bool) {
			// Append to (never overwrite) the current selection.
			for _, id := range ids {
				if id == "" || containsString(this.selectedSpells, id) {
					continue
				}

				this.selectedSpells = append(this.selectedSpells, id)
			}
			for len(this.spellRemoveBtns) < len(this.selectedSpells) {
				this.spellRemoveBtns = append(this.spellRemoveBtns, widget.Clickable{})
			}
		}))
	}

	if this.pickItemBtn.Clicked(gtx) {
		this.opener(NewItemPickerDialog("Pick Starting Item", nil, func(ids []string) {
			if len(ids) == 0 {
				return
			}

			if len(ids) == 1 {
				this.itemEdit.SetText(ids[0])
				return
			}

			receiver := this.receiver()
			entries := make([]config_inner.BonusEntry, 0, len(ids))
			for _, id := range ids {
				entries = append(entries, config_inner.BonusEntry{
					PresetType:     config_inner.BonusStartingItem,
					ReceiverFilter: receiver,
					Param:          id,
				})
			}
			this.pending = entries
			this.hasPending = true
		}))
	}
}

func (this *BonusPickerDialog) buildEntries() ([]config_inner.BonusEntry, bool) {
	this.errorText = ""
	presetType := this.getSelectedType()
	receiver := this.receiver()
	switch presetType {
	case config_inner.BonusTownPortalFree:
		return []config_inner.BonusEntry{{PresetType: presetType, ReceiverFilter: receiver}}, true
	case config_inner.BonusSpell:
		if len(this.selectedSpells) == 0 {
			this.errorText = "Pick at least one spell."
			return nil, false
		}
		param2 := "0"
		if this.makeFree.Value {
			param2 = "1"
		}
		entries := make([]config_inner.BonusEntry, 0, len(this.selectedSpells))
		for _, id := range this.selectedSpells {
			entries = append(entries,
				config_inner.BonusEntry{PresetType: presetType, ReceiverFilter: receiver, Param: id, Param2: param2})
		}
		return entries, true
	case config_inner.BonusUnitMultiplier:
		value := strings.TrimSpace(this.multiplierEdit.Text())
		if !isNumeric(value) {
			this.errorText = "Enter a numeric multiplier."
			return nil, false
		}
		return []config_inner.BonusEntry{{PresetType: presetType, ReceiverFilter: receiver, Param: value}}, true
	case config_inner.BonusMovementBonus:
		value := strings.TrimSpace(this.movementEdit.Text())
		if !isNumeric(value) {
			this.errorText = "Enter a numeric movement value."
			return nil, false
		}
		return []config_inner.BonusEntry{{PresetType: presetType, ReceiverFilter: receiver, Param: value}}, true
	case config_inner.BonusStartingItem:
		id := strings.TrimSpace(this.itemEdit.Text())
		if id == "" {
			this.errorText = "Pick or enter an item."
			return nil, false
		}
		return []config_inner.BonusEntry{{PresetType: presetType, ReceiverFilter: receiver, Param: id}}, true
	default:
		value := strings.TrimSpace(this.resourceEdit.Text())
		if !isNumeric(value) {
			this.errorText = "Enter a numeric amount."
			return nil, false
		}
		return []config_inner.BonusEntry{{PresetType: presetType, ReceiverFilter: receiver, Param: value}}, true
	}
}

func (this *BonusPickerDialog) getSelectedType() config_inner.BonusPresetType {
	index := this.typeDropdown.GetSelectedIndex()
	if index < 0 || index >= len(bonusTypeOptions) {
		return config_inner.BonusTownPortalFree
	}

	return bonusTypeOptions[index].typ
}

func (this *BonusPickerDialog) receiver() string {
	index := this.receiverDropdown.GetSelectedIndex()
	if index < 0 || index >= len(bonusReceiverOptions) {
		return bonusReceiverOptions[0]
	}

	return bonusReceiverOptions[index]
}

func (this *BonusPickerDialog) applyTypeDefaults(typ config_inner.BonusPresetType) {
	if def, ok := bonusResourceDefaults[typ]; ok {
		this.resourceEdit.SetText(def)
	}
}

// getSpellListWidget renders the picked spells as removable rows with a
// school-coloured bubble and the spell's display name.
func (this *BonusPickerDialog) getSpellListWidget(theme *material.Theme) layout.Widget {
	if len(this.selectedSpells) == 0 {
		return widgets.NewDimmedLabelWidget(theme, "(no spells picked - use Pick spells...)")
	}

	return func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.Y = min(gtx.Constraints.Max.Y, gtx.Dp(unit.Dp(150)))
		return material.List(theme, &this.spellScroll).Layout(gtx, len(this.selectedSpells),
			func(gtx layout.Context, index int) layout.Dimensions {
				return this.getSpellRowWidget(theme, index)(gtx)
			})
	}
}

func (this *BonusPickerDialog) getSpellRowWidget(theme *material.Theme, index int) layout.Widget {
	sid := this.selectedSpells[index]
	name, school := spellNameAndSchool(sid)
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.
			Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(constants.DefaultRoundnessLarge)
						paint.FillShape(gtx.Ops,
							constants.GetSpellSchoolColorFromDisplayName(school),
							clip.Ellipse{Max: image.Pt(size, size)}.Op(gtx.Ops))
						return layout.Dimensions{Size: image.Pt(size, size)}
					}),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
					layout.Rigid(widgets.NewLabelBuilder(theme).WithSizeBig().WithText(name).WithColorDefault().Build),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
					layout.Flexed(1, widgets.NewDimmedLabelWidget(theme, school)),
					layout.Rigid(widgets.NewButtonWidget(theme, "✕", &this.spellRemoveBtns[index], false)),
				)
			})
	}
}

// spellNameAndSchool resolves a spell SID to its display name and school
// label, with a sentence-case fallback for unknown SIDs.
func spellNameAndSchool(sid string) (name, school string) {
	if spell, ok := constants.FindSpell(sid); ok {
		label := constants.SpellSchoolDisplayNames[spell.School]
		if label == "" {
			label = spell.School
		}
		return spell.Name, label
	}
	return constants.SidToDisplayName(sid), "Spell"
}

func spellCountLabel(count int) string {
	switch count {
	case 0:
		return "Spells"
	case 1:
		return "1 spell picked"
	default:
		return strconv.Itoa(count) + " spells picked"
	}
}

func containsString(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

func isNumeric(value string) bool {
	if value == "" {
		return false
	}

	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
