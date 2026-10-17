package drivers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

const (
	conversionNoticeTitle = "Older file format converted"
	conversionNoticeBody  = "This file was saved in an older format. Values you did not edit in the zone editor " +
		"now follow the current settings and generator. Save to keep the new format."
)

func (this *State) ApplyEditedZones(request dtos.ZoneEditorZonesDto) {
	pendingBase := this.pendingBaseZones
	this.pendingBaseZones = dtos.ZoneEditorZonesDto{}
	if !this.hasTemplateVariants() {
		return
	}

	revertsToUntouchedBase := request.RevertToBase && matchesZoneSet(request, pendingBase)
	if !this.handleUpdateTemplate(request.Zones, request.Connections) {
		return
	}

	if revertsToUntouchedBase {
		this.clearDerivationBaseline()
		if this.innerState.ClearManualEdits() {
			this.flagAsUnsaved()
		}
		return
	}

	this.captureManualEdits()
}

func (this *State) PreviewBaseZones() (dtos.ZoneEditorZonesDto, bool) {
	dto, err := this.handler.GenerateTemplate(this.GetStateDto())
	if err != nil {
		this.SetStatus(fmt.Sprintf("Generation failed: %v.", err), true)
		return dtos.ZoneEditorZonesDto{}, false
	}

	if dto.Template == nil || len(dto.Template.Variants) == 0 {
		return dtos.ZoneEditorZonesDto{}, false
	}

	variant := dto.Template.Variants[0]
	this.pendingBaseZones = dtos.ZoneEditorZonesDto{
		Zones:       variant.Zones,
		Connections: variant.Connections,
	}
	return this.pendingBaseZones, true
}

func matchesZoneSet(left, right dtos.ZoneEditorZonesDto) bool {
	return reflect.DeepEqual(left.Zones, right.Zones) &&
		reflect.DeepEqual(left.Connections, right.Connections)
}

func (this *State) handleUpdateTemplate(zones []template_model.Zone, connections []template_model.Connection) bool {
	dto, err := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    this.lastTemplate,
		Zones:       zones,
		Connections: connections,
		EditorState: new(this.GetStateDto()),
	})

	if err != nil && errors.Is(err, common_errors.ErrProvidedTemplateInvalid) {
		this.SetStatus(
			fmt.Sprintf("Unable to update template, possibly because template was not generated. ‼ Error: %v", err),
			true)
		return false
	}

	this.setLastTemplate(dto.Template)
	if err != nil {
		this.SetStatus(
			fmt.Sprintf(
				"Applied %d zones and %d connections. ‼ Error: %v; fix before export.",
				len(zones), len(connections), err),
			true)
		return true
	}

	this.SetStatus(
		fmt.Sprintf("Applied %d zones and %d connections from the editor.", len(zones), len(connections)),
		false)
	return true
}

// captureManualEdits records the applied layout against a fresh generation. On
// failure the layout is kept whole (all-overridden) and converted later.
func (this *State) captureManualEdits() {
	variant := this.lastTemplate.Variants[0]
	settings := this.innerState.GetCurrentSettings()
	captured, err := this.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{
		State:       editor_state_dto.EditorStateDto{EditorState: settings},
		Zones:       variant.Zones,
		Connections: variant.Connections,
	})
	if err != nil {
		this.innerState.SetManualEdits(
			editor_state_model.ManualEditSettingsFromSnapshot(variant.Zones, variant.Connections))
		this.innerState.SetGeneratorVersion(0)
		this.derivationBaseline = &settings
		this.flagAsUnsaved()
		this.SetStatus(
			fmt.Sprintf("Applied %d zones and %d connections. ‼ Error: %v",
				len(variant.Zones), len(variant.Connections), err),
			true)
		return
	}

	this.clearDerivationBaseline()
	this.innerState.SetGeneratorVersion(common_generator.GeneratorOutputVersion)
	if this.innerState.SetManualEdits(captured.ManualEdits) {
		this.flagAsUnsaved()
	}
}

// reapplyManualEdits rebuilds the manual edits on the freshly generated
// template, converting a pending snapshot first. It returns the notes for the
// generation status and whether a conversion failed.
func (this *State) reapplyManualEdits(decision dtos.ManualEditDecisionDto) ([]string, bool) {
	castleChanges := *decision.ReapplyWithCastleChanges
	if this.innerState.NeedsDerivation() {
		baselineChanges, err := this.convertManualEdits()
		if err != nil {
			zones, connections := this.innerState.GetManualEdits().ToSnapshot()
			this.handleUpdateTemplate(zones, connections)
			return []string{fmt.Sprintf("Manual edits could not be converted: %v.", err)}, true
		}
		castleChanges = baselineChanges
	}

	expansion := this.handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{
		Template:      this.lastTemplate,
		State:         this.GetStateDto(),
		CastleChanges: castleChanges,
	})
	this.handleUpdateTemplate(expansion.Zones, expansion.Connections)
	if expansion.Changed && this.innerState.SetManualEdits(expansion.ManualEdits) {
		this.flagAsUnsaved()
	}

	notes := expansion.Warnings
	if decision.GeneratorVersionChanged {
		notes = append(notes, fmt.Sprintf(
			"Saved with generator version %d; values you did not edit now follow version %d and may differ.",
			this.innerState.GetGeneratorVersion(), common_generator.GeneratorOutputVersion))
		this.innerState.SetGeneratorVersion(common_generator.GeneratorOutputVersion)
	}
	return notes, false
}

// convertManualEdits records the pending snapshot against its baseline
// settings, or the current ones when the baseline cannot be generated, and
// returns the castle options that moved since that baseline.
func (this *State) convertManualEdits() (editor_state_model.CastleSettingChanges, error) {
	current := this.innerState.GetCurrentSettings()
	baseline := current
	if this.derivationBaseline != nil {
		baseline = this.derivationBaseline.CloneWithoutManualEdits()
	}

	edits := this.innerState.GetManualEdits()
	converted, err := this.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{
		Baseline:    editor_state_dto.EditorStateDto{EditorState: baseline},
		ManualEdits: edits,
	})
	if err != nil && this.derivationBaseline != nil {
		baseline = current
		converted, err = this.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{
			Baseline:    editor_state_dto.EditorStateDto{EditorState: current},
			ManualEdits: edits,
		})
	}
	if err != nil {
		return editor_state_model.CastleSettingChanges{}, err
	}

	this.innerState.SetManualEdits(converted.ManualEdits)
	this.innerState.SetGeneratorVersion(common_generator.GeneratorOutputVersion)
	this.derivationBaseline = nil
	this.flagAsUnsaved()
	if this.announceConversion {
		this.announceConversion = false
		this.dialogs.Open(dialogs.NewMessageDialog(conversionNoticeTitle, conversionNoticeBody))
	}
	return baseline.DiffCastleSettings(&current), nil
}
