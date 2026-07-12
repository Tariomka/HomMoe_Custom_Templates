package dialogs

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
)

func (this *ZoneEditorDialog) propertyRows(theme *material.Theme) []layout.Widget {
	connection := this.selected
	rows := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, connection.From+"  →  "+connection.To)
			label.Color = themes.ColorAccentBright
			label.Font = font.Font{Weight: font.SemiBold}
			return label.Layout(gtx)
		},
	}
	if connection.IsUserAdded {
		rows = append(rows, widgets.NewDimmedLabelWidget(theme, "User-added connection"))
	}
	rows = append(
		rows,
		widgets.NewVerticalSpacerWidget(6),
		widgets.NewLabeledRowWidget(theme, "Type", 110, this.typeDropdown.GetWidget(theme)),
		widgets.NewLabeledRowWidget(theme, "Guard zone", 110, this.guardZoneDropdown.GetWidget(theme)),
		widgets.NewVerticalSpacerWidget(4),
		widgets.NewLabeledRowWidget(theme, "Guard preset", 110, this.guardDropdown.GetWidget(theme)),
		widgets.NewLabeledRowWidget(
			theme,
			"Guard value",
			110,
			widgets.NewTextboxWidget(theme, &this.guardValueEdit, "guard value", false),
		),
		widgets.NewVerticalSpacerWidget(4),
		widgets.NewLabeledRowWidget(theme, "Weekly +", 110, this.weeklyDropdown.GetWidget(theme)),
		widgets.NewLabeledRowWidget(
			theme,
			"Increment",
			110,
			widgets.NewTextboxWidget(theme, &this.weeklyEdit, "0.15", false),
		),
		widgets.NewVerticalSpacerWidget(6),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.advancedBool, "Advanced options"),
	)
	if this.advancedBool.Value {
		rows = append(
			rows,
			widgets.NewLabeledRowWidget(
				theme,
				"Match group",
				110,
				widgets.NewTextboxWidget(theme, &this.matchGroupEdit, "rnd_guard_...", false),
			),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.escapeBool, "Guard escape"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.simSquadBool, "Sim turn squad"),
		)
	}
	rows = append(rows,
		widgets.NewVerticalSpacerWidget(10),
		widgets.NewButtonWidget(theme, "Delete this connection", &this.sidePropDelete, false),
	)
	return rows
}

// syncPropsFromConnection loads the property widgets from the selected connection.
// Called once whenever the selection changes.
func (this *ZoneEditorDialog) syncPropsFromConnection() {
	connection := this.selected
	if connection == nil {
		return
	}
	this.typeDropdown.SetItems(connection_editor.UserCreatableConnectionTypes())
	if !this.typeDropdown.SelectByName(connection.ConnectionType) {
		this.typeDropdown.SelectByName("Direct")
	}
	this.guardZoneDropdown.SetItems([]string{connection.From, connection.To})
	if !this.guardZoneDropdown.SelectByName(connection.GuardZone) {
		this.guardZoneDropdown.SelectByName(connection.From)
	}
	tier := connection_editor.HigherTierOf(connection.From, connection.To, this.zones, this.playerZones)
	labels, values := guardPresetItems(tier)
	this.guardPresetValues = values
	this.guardDropdown.SetItems(labels)
	this.guardDropdown.SelectByName(matchGuardLabel(labels, values, connection.GuardValue))
	this.guardValueEdit.SetText(strconv.Itoa(connection.GuardValue))
	this.weeklyDropdown.SetItems(connection_editor.WeeklyIncrementLabels)
	this.weeklyDropdown.SelectByName(matchWeeklyLabel(connection.GuardWeeklyIncrement))
	this.weeklyEdit.SetText(formatIncrement(connection.GuardWeeklyIncrement))
	this.matchGroupEdit.SetText(connection.GuardMatchGroup)
	this.escapeBool.Value = connection.GuardEscape
	this.simSquadBool.Value = connection.SimTurnSquad
}

// writebackProps copies the property widget state back into the selected
// connection after the panel has been laid out for this frame.
func (this *ZoneEditorDialog) writebackProps() {
	connection := this.selected
	if connection == nil {
		return
	}
	typeItems := connection_editor.UserCreatableConnectionTypes()
	if index := this.typeDropdown.GetSelectedIndex(); index >= 0 && index < len(typeItems) {
		connection.ConnectionType = typeItems[index]
	}
	zoneItems := []string{connection.From, connection.To}
	if index := this.guardZoneDropdown.GetSelectedIndex(); index >= 0 && index < len(zoneItems) {
		connection.GuardZone = zoneItems[index]
	}
	if this.guardDropdown.WasUpdated {
		if index := this.guardDropdown.GetSelectedIndex(); index >= 0 && index < len(this.guardPresetValues) {
			this.guardValueEdit.SetText(strconv.Itoa(this.guardPresetValues[index]))
		}
	}
	if value, err := strconv.Atoi(strings.TrimSpace(this.guardValueEdit.Text())); err == nil {
		connection.GuardValue = value
	}
	if this.weeklyDropdown.WasUpdated {
		if index := this.weeklyDropdown.GetSelectedIndex(); index >= 0 &&
			index < len(connection_editor.WeeklyIncrementValues) {
			this.weeklyEdit.SetText(formatIncrement(connection_editor.WeeklyIncrementValues[index]))
		}
	}
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.weeklyEdit.Text()), 64); err == nil {
		connection.GuardWeeklyIncrement = value
	}
	connection.GuardMatchGroup = strings.TrimSpace(this.matchGroupEdit.Text())
	connection.GuardEscape = this.escapeBool.Value
	connection.SimTurnSquad = this.simSquadBool.Value
}

func guardPresetItems(tier connection_editor.ZoneTier) (labels []string, values []int) {
	for _, extra := range connection_editor.ExtrasForTier(tier) {
		labels = append(labels, fmt.Sprintf("%s (%d)", extra.Label, extra.Value))
		values = append(values, extra.Value)
	}
	presets := connection_editor.GuardPresetsForTier(tier)
	for i, strength := range connection_editor.StrengthLabels {
		labels = append(labels, fmt.Sprintf("%s (%d)", strength, presets[i]))
		values = append(values, presets[i])
	}
	return labels, values
}

func matchGuardLabel(labels []string, values []int, value int) string {
	for i, candidate := range values {
		if candidate == value {
			return labels[i]
		}
	}
	return ""
}

func matchWeeklyLabel(value float64) string {
	for i, candidate := range connection_editor.WeeklyIncrementValues {
		if math.Abs(candidate-value) < 1e-9 {
			return connection_editor.WeeklyIncrementLabels[i]
		}
	}
	return ""
}

func formatIncrement(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}
