package dialogs

import (
	"image"
	"slices"

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
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/config_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// BonusPickerDialog composes one or more game-start bonuses. For the Spell
// type it keeps a removable list of picked spells (the spell picker appends to
// it) and emits one bonus per spell; for StartingItem it can launch the item
// picker (via opener). Implements widgets.Dialog.
type BonusPickerDialog struct {
	existingKeys     map[string]bool
	existingSpellIDs []string
	opener           interfaces.DialogOpener
	handler          handler_interfaces.IGuiHandler
	onApply          func(entries []config.BonusEntry)

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
	pending    []config.BonusEntry
	hasPending bool
}

// NewBonusPickerDialog builds the bonus composer. existing is used for
// duplicate detection and to pre-exclude already-chosen spells in the spell
// sub-picker. onApply receives the composed bonus entries.
func NewBonusPickerDialog(
	existing []config.BonusEntry,
	opener interfaces.DialogOpener,
	handler handler_interfaces.IGuiHandler,
	onApply func(entries []config.BonusEntry)) *BonusPickerDialog {
	summary := handler.DescribeExistingBonuses(existing)

	labels := linq.FromSlice(constants.GetBonusTypeOptions()).
		Select(func(opt constants.BonusTypeOption) string { return opt.Label }).
		ToSlice()

	dialog := &BonusPickerDialog{
		existingKeys:     summary.Keys,
		existingSpellIDs: summary.SpellIDs,
		opener:           opener,
		handler:          handler,
		onApply:          onApply,
		typeDropdown:     components.NewDropdownSelector(labels),
		receiverDropdown: components.NewDropdownSelector(append([]string{}, constants.GetBonusReceiverOptions()...)),
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

func (this *BonusPickerDialog) PreferredSize() (unit.Dp, unit.Dp) { return unit.Dp(460), unit.Dp(420) }

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
				this.applyTypeDefaults(this.getSelectedType())
			}
			return dims
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if config_helpers.IsResource(presetType) {
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
				widgets.NewLabelWidget(theme, this.errorText, themes.ColorsBase.Error))
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

func (this *BonusPickerDialog) getEditorWidget(
	theme *material.Theme,
	presetType config.BonusPresetType) layout.Widget {
	switch presetType {
	case config.BonusTownPortalFree:
		return widgets.NewDimmedLabelWidget(
			theme,
			"Grants free town portal to the selected hero(es). No extra parameters.",
		)
	case config.BonusSpell:
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, widgets.NewDimmedLabelWidget(theme,
							this.handler.GetSpellCountLabel(len(this.selectedSpells)))),
						widgets.NewDefaultComponentSpacer(),
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
	case config.BonusUnitMultiplier:
		return widgets.NewLabeledRowWidget(theme, "Multiplier", 90,
			widgets.NewTextboxWidget(theme, &this.multiplierEdit, "2", false))
	case config.BonusMovementBonus:
		return widgets.NewLabeledRowWidget(theme, "Movement", 90,
			widgets.NewTextboxWidget(theme, &this.movementEdit, "0", false))
	case config.BonusStartingItem:
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.itemEdit, "use Pick item...", false)),
				widgets.NewDefaultComponentSpacer(),
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

	if !this.addBtn.Clicked(gtx) {
		return false
	}

	result := this.handler.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     this.getSelectedType(),
		ReceiverFilter: this.receiver(),
		SelectedSpells: this.selectedSpells,
		MakeSpellsFree: this.makeFree.Value,
		MultiplierText: this.multiplierEdit.Text(),
		MovementText:   this.movementEdit.Text(),
		ItemText:       this.itemEdit.Text(),
		ResourceText:   this.resourceEdit.Text(),
	})
	this.errorText = result.Error
	if result.Error != "" {
		return false
	}

	fresh := this.handler.FilterNewBonusEntries(result.Entries, this.existingKeys)
	if len(fresh) == 0 {
		this.errorText = "That bonus already exists."
		return false
	}

	if this.onApply != nil {
		this.onApply(fresh)
	}
	return true
}

// handleSubPickers reacts to the "Pick spell" / "Pick item" buttons by pushing
// the relevant picker onto the dialog stack.
func (this *BonusPickerDialog) handleSubPickers(gtx layout.Context) {
	if this.opener == nil {
		return
	}

	if this.pickSpellBtn.Clicked(gtx) {
		excluded := append(append([]string{}, this.existingSpellIDs...), this.selectedSpells...)
		this.opener(NewSpellPickerDialog(excluded, false, this.handler, func(ids []string, _ bool) {
			// Append to (never overwrite) the current selection.
			for _, id := range ids {
				if id != "" && !slices.Contains(this.selectedSpells, id) {
					this.selectedSpells = append(this.selectedSpells, id)
				}
			}
			for len(this.spellRemoveBtns) < len(this.selectedSpells) {
				this.spellRemoveBtns = append(this.spellRemoveBtns, widget.Clickable{})
			}
		}))
	}

	if this.pickItemBtn.Clicked(gtx) {
		this.opener(NewItemPickerDialog("Pick Starting Item", nil, this.handler, func(ids []string) {
			if len(ids) == 0 {
				return
			}

			if len(ids) == 1 {
				this.itemEdit.SetText(ids[0])
				return
			}

			receiver := this.receiver()
			entries := make([]config.BonusEntry, 0, len(ids))
			for _, id := range ids {
				entries = append(entries, config.BonusEntry{
					PresetType:     config.BonusStartingItem,
					ReceiverFilter: receiver,
					Param:          id,
				})
			}
			this.pending = entries
			this.hasPending = true
		}))
	}
}

func (this *BonusPickerDialog) getSelectedType() config.BonusPresetType {
	index := this.typeDropdown.GetSelectedIndex()
	options := constants.GetBonusTypeOptions()
	if index < 0 || index >= len(options) {
		return config.BonusTownPortalFree
	}

	return options[index].PresetType
}

func (this *BonusPickerDialog) receiver() string {
	index := this.receiverDropdown.GetSelectedIndex()
	options := constants.GetBonusReceiverOptions()
	if index < 0 || index >= len(options) {
		return options[0]
	}

	return options[index]
}

func (this *BonusPickerDialog) applyTypeDefaults(typ config.BonusPresetType) {
	if def, ok := constants.GetBonusResourceDefaults()[typ]; ok {
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
	name, school := constants.GetSpellNameAndSchool(sid)
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
					widgets.NewDefaultComponentSpacer(),
					layout.Rigid(widgets.NewLabelBuilder(theme).WithSizeBig().WithText(name).WithColorDefault().Build),
					widgets.NewDefaultComponentSpacer(),
					layout.Flexed(1, widgets.NewDimmedLabelWidget(theme, school)),
					layout.Rigid(widgets.NewButtonWidget(theme, "X", &this.spellRemoveBtns[index], false)),
				)
			})
	}
}
