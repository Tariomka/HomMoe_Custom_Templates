package drivers

import (
	"errors"
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
)

// ApplyEditedZones writes zones and connections edited in the manual zone
// editor back into the live template and stores them in the editor state as
// the authoritative manual snapshot, reapplied on later regenerations and
// saved with the rest of the .gen.json state.
func (this *State) ApplyEditedZones(zones []entities.Zone, connections []entities.Connection) {
	if !this.hasTemplateVariants() {
		return
	}

	this.handleUpdateTemplate(zones, connections)
	this.innerState.SetManualEdits(zones, connections)
}

func (this *State) handleUpdateTemplate(zones []entities.Zone, connections []entities.Connection) {
	dto, err := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    this.lastTemplate,
		Zones:       zones,
		Connections: connections,
		Config:      this.GetGeneratorConfig(),
	})

	if err != nil && errors.Is(err, common.ErrProvidedTemplateInvalid) {
		this.SetStatus(
			fmt.Sprintf("Unable to update template, possibly because template was not generated. ⚠ Error: %v", err),
			true)
		return
	}

	this.lastTemplate = dto.Template
	if err != nil {
		this.SetStatus(
			fmt.Sprintf(
				"Applied %d zones and %d connections. ⚠ Error: %v; fix before export.",
				len(zones), len(connections), err),
			true)
		return
	}

	this.SetStatus(
		fmt.Sprintf("Applied %d zones and %d connections from the editor.", len(zones), len(connections)),
		false)
}

// reapplyManualEdits restores the manual zone/connection snapshot over the
// freshly generated template. When castle-count options changed since the
// last generation - the only generator options that override manual edits -
// the new counts are first pushed into the snapshot and the updated snapshot
// is stored back so later regenerations and saves carry it.
func (this *State) reapplyManualEdits(castleChanges editor_state_dto.CastleSettingChanges) {
	zones := this.innerState.GetManualZones()
	connections := this.innerState.GetManualConnections()
	if castleChanges.Any() {
		connection_editor.ApplyCastleSettingChanges(zones, castleChanges, this.GetGeneratorConfig())
		this.innerState.SetManualEdits(zones, connections)
	}
	this.handleUpdateTemplate(zones, connections)
}
