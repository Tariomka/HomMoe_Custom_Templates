package components

import (
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
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

// BonusPickerDialog composes a single game-start bonus. For Spell and
// StartingItem types it can launch the spell / item pickers (via opener); when
// those return more than one selection it emits one bonus per id and closes.
// Implements widgets.Dialog.
type BonusPickerDialog struct {
	existingKeys     map[string]bool
	existingSpellIds []string
	opener           widgets.DialogOpener
	onApply          func(entries []config_inner.BonusEntry)

	typeDropdown     *content.DropdownSelector
	receiverDropdown *content.DropdownSelector

	spellEdit    widget.Editor
	makeFree     widget.Bool
	pickSpellBtn widget.Clickable

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
func NewBonusPickerDialog(existing []config_inner.BonusEntry, opener widgets.DialogOpener, onApply func(entries []config_inner.BonusEntry)) *BonusPickerDialog {
	keys := make(map[string]bool, len(existing))
	var spellIds []string
	for _, entry := range existing {
		keys[entry.String()] = true
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
		typeDropdown:     content.NewDropdownSelector(labels),
		receiverDropdown: content.NewDropdownSelector(append([]string{}, bonusReceiverOptions...)),
	}
	dialog.spellEdit.SingleLine = true
	dialog.spellEdit.ReadOnly = true
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
	// A multi-selection sub-picker queued several bonuses; emit and close.
	if this.hasPending {
		this.hasPending = false
		if this.onApply != nil {
			this.onApply(this.pending)
		}
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}
	if this.cancelBtn.Clicked(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}
	if this.addBtn.Clicked(gtx) {
		if entry, ok := this.buildEntry(); ok {
			if this.existingKeys[entry.String()] {
				this.errorText = "That bonus already exists."
			} else {
				if this.onApply != nil {
					this.onApply([]config_inner.BonusEntry{entry})
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}, true
			}
		}
	}
	this.handleSubPickers(gtx)

	selectedType := this.selectedType()

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := widgets.NewLabeledRowWidget(theme, "Type", 90, func(gtx layout.Context) layout.Dimensions {
				return this.typeDropdown.Layout(gtx, theme)
			})(gtx)
			if this.typeDropdown.WasUpdated {
				this.applyTypeDefaults(this.selectedType())
			}
			return dims
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if isResourceBonus(this.selectedType()) {
				return layout.Dimensions{}
			}
			return widgets.NewLabeledRowWidget(theme, "Receiver", 90, func(gtx layout.Context) layout.Dimensions {
				return this.receiverDropdown.Layout(gtx, theme)
			})(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return this.layoutEditor(gtx, theme, selectedType)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if this.errorText == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, this.errorText)
				label.Color = themes.ColorError
				label.TextSize = unit.Sp(12)
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: gtx.Constraints.Min}
				}),
				layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelBtn, false)),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
				layout.Rigid(widgets.NewGoldButtonWidget(theme, "Add", &this.addBtn, false)),
			)
		}),
	), false
}

func (this *BonusPickerDialog) layoutEditor(gtx layout.Context, theme *material.Theme, typ config_inner.BonusPresetType) layout.Dimensions {
	switch typ {
	case config_inner.BonusTownPortalFree:
		label := material.Body2(theme, "Grants free town portal to the selected hero(es). No extra parameters.")
		label.Color = themes.ColorTextDim
		label.TextSize = unit.Sp(12)
		return label.Layout(gtx)
	case config_inner.BonusSpell:
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.spellEdit, "use Pick spell…")),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
					layout.Rigid(widgets.NewButtonWidget(theme, "Pick spell", &this.pickSpellBtn, this.opener == nil)),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.makeFree, "Make spell free")),
		)
	case config_inner.BonusUnitMultiplier:
		return widgets.NewLabeledRowWidget(theme, "Multiplier", 90, widgets.NewTextboxWidget(theme, &this.multiplierEdit, "2"))(gtx)
	case config_inner.BonusMovementBonus:
		return widgets.NewLabeledRowWidget(theme, "Movement", 90, widgets.NewTextboxWidget(theme, &this.movementEdit, "0"))(gtx)
	case config_inner.BonusStartingItem:
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.itemEdit, "use Pick item…")),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Pick item", &this.pickItemBtn, this.opener == nil)),
		)
	default: // resources
		return widgets.NewLabeledRowWidget(theme, "Amount", 90, widgets.NewTextboxWidget(theme, &this.resourceEdit, "0"))(gtx)
	}
}

// handleSubPickers reacts to the "Pick spell" / "Pick item" buttons by pushing
// the relevant picker onto the dialog stack.
func (this *BonusPickerDialog) handleSubPickers(gtx layout.Context) {
	if this.opener == nil {
		return
	}
	if this.pickSpellBtn.Clicked(gtx) {
		this.opener(NewSpellPickerDialog(this.existingSpellIds, true, func(ids []string, makeFree bool) {
			if len(ids) == 0 {
				return
			}
			if len(ids) == 1 {
				this.spellEdit.SetText(ids[0])
				this.makeFree.Value = makeFree
				return
			}
			param2 := "0"
			if makeFree {
				param2 = "1"
			}
			receiver := this.receiver()
			entries := make([]config_inner.BonusEntry, 0, len(ids))
			for _, id := range ids {
				entries = append(entries, config_inner.BonusEntry{
					PresetType:     config_inner.BonusSpell,
					ReceiverFilter: receiver,
					Param:          id,
					Param2:         param2,
				})
			}
			this.pending = entries
			this.hasPending = true
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

func (this *BonusPickerDialog) buildEntry() (config_inner.BonusEntry, bool) {
	this.errorText = ""
	typ := this.selectedType()
	receiver := this.receiver()
	switch typ {
	case config_inner.BonusTownPortalFree:
		return config_inner.BonusEntry{PresetType: typ, ReceiverFilter: receiver}, true
	case config_inner.BonusSpell:
		id := strings.TrimSpace(this.spellEdit.Text())
		if id == "" {
			this.errorText = "Pick or enter a spell."
			return config_inner.BonusEntry{}, false
		}
		param2 := "0"
		if this.makeFree.Value {
			param2 = "1"
		}
		return config_inner.BonusEntry{PresetType: typ, ReceiverFilter: receiver, Param: id, Param2: param2}, true
	case config_inner.BonusUnitMultiplier:
		value := strings.TrimSpace(this.multiplierEdit.Text())
		if !isNumeric(value) {
			this.errorText = "Enter a numeric multiplier."
			return config_inner.BonusEntry{}, false
		}
		return config_inner.BonusEntry{PresetType: typ, ReceiverFilter: receiver, Param: value}, true
	case config_inner.BonusMovementBonus:
		value := strings.TrimSpace(this.movementEdit.Text())
		if !isNumeric(value) {
			this.errorText = "Enter a numeric movement value."
			return config_inner.BonusEntry{}, false
		}
		return config_inner.BonusEntry{PresetType: typ, ReceiverFilter: receiver, Param: value}, true
	case config_inner.BonusStartingItem:
		id := strings.TrimSpace(this.itemEdit.Text())
		if id == "" {
			this.errorText = "Pick or enter an item."
			return config_inner.BonusEntry{}, false
		}
		return config_inner.BonusEntry{PresetType: typ, ReceiverFilter: receiver, Param: id}, true
	default:
		value := strings.TrimSpace(this.resourceEdit.Text())
		if !isNumeric(value) {
			this.errorText = "Enter a numeric amount."
			return config_inner.BonusEntry{}, false
		}
		return config_inner.BonusEntry{PresetType: typ, ReceiverFilter: receiver, Param: value}, true
	}
}

func (this *BonusPickerDialog) selectedType() config_inner.BonusPresetType {
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

func isResourceBonus(typ config_inner.BonusPresetType) bool {
	return typ >= config_inner.BonusStartingGold
}

func isNumeric(value string) bool {
	if value == "" {
		return false
	}
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
