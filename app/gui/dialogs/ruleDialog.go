package dialogs

import (
	"image"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// ManageRulesDialog edits the polymorphic content-rule list for a single
// zone-content row. It mirrors the C# Manage/Add zone-content-rule windows,
// folded into one modal: the upper area lists applied rules, the lower area
// adds or updates a rule (one per type).
type ManageRulesDialog struct {
	mapping models.SidMapping
	rules   []models.ContentRuleRowSave
	onApply func([]models.ContentRuleRowSave)

	contentRuleHandler handler_interfaces.IContentRuleHandler
	types              []dtos.ContentRuleOptionDto
	distanceNames      []string
	variantIDs         []int
	variantLabels      []string

	scroll           widget.List
	typeDropdown     *components.DropdownSelector
	distanceDropdown *components.DropdownSelector
	variantDropdown  *components.DropdownSelector
	guardedBool      widget.Bool
	soloBool         widget.Bool

	addBtn     widget.Clickable
	applyBtn   widget.Clickable
	cancelBtn  widget.Clickable
	removeBtns []widget.Clickable
}

// NewManageRulesDialog builds the dialog for the given content row. onApply is
// invoked with the edited rule list when the user clicks Apply.
func NewManageRulesDialog(
	mapping models.SidMapping,
	rules []models.ContentRuleRowSave,
	contentRuleHandler handler_interfaces.IContentRuleHandler,
	onApply func([]models.ContentRuleRowSave)) *ManageRulesDialog {
	options := contentRuleHandler.GetContentRuleEditorOptions(mapping)
	dialog := &ManageRulesDialog{
		mapping:            mapping,
		rules:              utils.CloneRuleRows(rules),
		onApply:            onApply,
		contentRuleHandler: contentRuleHandler,
		types:              options.Rules,
		distanceNames:      options.Distances,
	}
	dialog.scroll.Axis = layout.Vertical

	for _, variant := range options.Variants {
		dialog.variantIDs = append(dialog.variantIDs, variant.ID)
		dialog.variantLabels = append(dialog.variantLabels, variant.Label)
	}

	typeLabels := make([]string, len(dialog.types))
	for i, option := range dialog.types {
		typeLabels[i] = option.Name
	}
	dialog.typeDropdown = components.NewDropdownSelector(typeLabels)
	dialog.distanceDropdown = components.NewDropdownSelector(dialog.distanceNames)
	dialog.variantDropdown = components.NewDropdownSelector(dialog.variantLabels)
	dialog.guardedBool.Value = true
	return dialog
}

func (this *ManageRulesDialog) Title() string {
	return "Content Rules - " + this.mapping.Name
}

func (this *ManageRulesDialog) PreferredSize() (unit.Dp, unit.Dp) {
	return unit.Dp(540), unit.Dp(500)
}

func (this *ManageRulesDialog) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	if this.applyBtn.Clicked(gtx) {
		if this.onApply != nil {
			this.onApply(utils.CloneRuleRows(this.rules))
		}
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}
	if this.cancelBtn.Clicked(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}
	if this.addBtn.Clicked(gtx) {
		this.upsertFromEditor()
	}
	this.processRemovals(gtx)

	rows := this.buildContentWidgets(theme)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(theme, &this.scroll).Layout(gtx, len(rows),
				func(gtx layout.Context, index int) layout.Dimensions {
					return rows[index](gtx)
				})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, 0)}
				}),
				layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelBtn, false)),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
				layout.Rigid(widgets.NewBrightButtonWidget(theme, "Apply", &this.applyBtn, false)),
			)
		}),
	), false
}

func (this *ManageRulesDialog) processRemovals(gtx layout.Context) {
	for len(this.removeBtns) < len(this.rules) {
		this.removeBtns = append(this.removeBtns, widget.Clickable{})
	}
	for i := range this.rules {
		if i < len(this.removeBtns) && this.removeBtns[i].Clicked(gtx) {
			this.rules = append(this.rules[:i:i], this.rules[i+1:]...)
			return
		}
	}
}

func (this *ManageRulesDialog) buildContentWidgets(theme *material.Theme) []layout.Widget {
	rows := []layout.Widget{this.sectionLabel(theme, "Applied rules")}

	if len(this.rules) == 0 {
		rows = append(rows, this.dimLabel(theme, "(no rules - add one below)"))
	} else {
		for i := range this.rules {
			rows = append(rows, this.layoutRuleRow(theme, i))
		}
	}

	rows = append(rows,
		func(gtx layout.Context) layout.Dimensions { return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx) },
		this.sectionLabel(theme, "Add / update rule"),
		this.layoutTypeRow(theme),
		this.layoutEditor(theme),
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(6)}.Layout(gtx,
				widgets.NewButtonWidget(theme, "+ Add / update", &this.addBtn, false))
		},
	)
	return rows
}

func (this *ManageRulesDialog) layoutRuleRow(theme *material.Theme, index int) layout.Widget {
	saved := this.rules[index]
	return func(gtx layout.Context) layout.Dimensions {
		return widgets.NewPanelWidget(unit.Dp(6), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(theme, this.ruleDisplayText(saved))
					label.Color = themes.ColorsBase.Text
					label.TextSize = unit.Sp(13)
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if index < len(this.removeBtns) {
						return widgets.NewButtonWidget(theme, "Remove", &this.removeBtns[index], false)(gtx)
					}
					return layout.Dimensions{}
				}),
			)
		})(gtx)
	}
}

func (this *ManageRulesDialog) layoutTypeRow(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(90)
				return this.dimLabel(theme, "Rule type")(gtx)
			}),
			layout.Flexed(1, this.typeDropdown.GetWidget(theme)),
		)
	}
}

func (this *ManageRulesDialog) layoutEditor(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		option, ok := this.selectedRuleType()
		if !ok {
			return layout.Dimensions{}
		}
		return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			switch option.EditorKind {
			case dtos.ContentRuleEditorKindDistance:
				return this.labeledControl(theme, option.EditorLabel, this.distanceDropdown.GetWidget(theme))(gtx)
			case dtos.ContentRuleEditorKindBoolean:
				if option.Key == dtos.ContentRuleKeyGuarded {
					return widgets.NewLabeledCheckboxRowWidget(theme, &this.guardedBool, option.EditorLabel)(gtx)
				}
				return widgets.NewLabeledCheckboxRowWidget(theme, &this.soloBool, option.EditorLabel)(gtx)
			case dtos.ContentRuleEditorKindVariant:
				return this.labeledControl(theme, option.EditorLabel, this.variantDropdown.GetWidget(theme))(gtx)
			default:
				return layout.Dimensions{}
			}
		})
	}
}

func (this *ManageRulesDialog) labeledControl(
	theme *material.Theme,
	label string,
	control layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(90)
				return this.dimLabel(theme, label)(gtx)
			}),
			layout.Flexed(1, control),
		)
	}
}

func (this *ManageRulesDialog) sectionLabel(theme *material.Theme, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(theme, text)
				label.Color = themes.ColorsBase.Accent
				label.TextSize = unit.Sp(13)
				return label.Layout(gtx)
			})
	}
}

func (this *ManageRulesDialog) dimLabel(theme *material.Theme, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		label := material.Body2(theme, text)
		label.Color = themes.ColorsBase.TextDim
		label.TextSize = unit.Sp(12)
		return label.Layout(gtx)
	}
}

// upsertFromEditor builds a rule from the editor state and inserts or replaces
// the rule of the same type, enforcing one-rule-per-type.
func (this *ManageRulesDialog) upsertFromEditor() {
	saved, ok := this.buildRuleFromEditor()
	if !ok {
		return
	}
	for i := range this.rules {
		if strings.EqualFold(this.rules[i].Name, saved.Name) {
			this.rules[i] = saved
			return
		}
	}
	this.rules = append(this.rules, saved)
}

func (this *ManageRulesDialog) buildRuleFromEditor() (models.ContentRuleRowSave, bool) {
	option, ok := this.selectedRuleType()
	if !ok {
		return models.ContentRuleRowSave{}, false
	}
	switch option.Key {
	case dtos.ContentRuleKeyDistanceToRoad, dtos.ContentRuleKeyDistanceToTown:
		distIdx := this.distanceDropdown.GetSelectedIndex()
		if distIdx < 0 || distIdx >= len(this.distanceNames) {
			return models.ContentRuleRowSave{}, false
		}
		return models.ContentRuleRowSave{Name: option.Name, DistanceName: this.distanceNames[distIdx]}, true
	case dtos.ContentRuleKeyGuarded:
		value := this.guardedBool.Value
		return models.ContentRuleRowSave{Name: option.Name, IsGuarded: &value}, true
	case dtos.ContentRuleKeySoloEncounter:
		value := this.soloBool.Value
		return models.ContentRuleRowSave{Name: option.Name, IsSoloEncounter: &value}, true
	case dtos.ContentRuleKeyVariant:
		variantIdx := this.variantDropdown.GetSelectedIndex()
		if variantIdx < 0 || variantIdx >= len(this.variantIDs) {
			return models.ContentRuleRowSave{}, false
		}
		id := this.variantIDs[variantIdx]
		return models.ContentRuleRowSave{Name: option.Name, VariantID: &id}, true
	}
	return models.ContentRuleRowSave{}, false
}

func (this *ManageRulesDialog) selectedRuleType() (dtos.ContentRuleOptionDto, bool) {
	index := this.typeDropdown.GetSelectedIndex()
	if index < 0 || index >= len(this.types) {
		return dtos.ContentRuleOptionDto{}, false
	}
	return this.types[index], true
}

// ruleDisplayText reconstructs a rule's user-facing description, falling back to
// the raw name when the saved data cannot be resolved to a known rule.
func (this *ManageRulesDialog) ruleDisplayText(saved models.ContentRuleRowSave) string {
	return this.contentRuleHandler.DescribeContentRule(this.mapping, saved).DisplayText
}
