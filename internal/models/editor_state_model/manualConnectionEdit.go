package editor_state_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ManualConnectionEdit is one manually edited connection in the C3 form: its
// structure (name, endpoints, type, provenance) plus the fields that differ
// from the connection it derives from.
type ManualConnectionEdit struct {
	editor_state.ManualConnectionEdit
}

// NewManualConnectionEdit records edited as structure plus its differences from derived.
func NewManualConnectionEdit(derived, edited template_model.Connection) ManualConnectionEdit {
	overrides, unset := DiffFieldOverrides[template_entity.Connection, editor_state.ConnectionFieldOverrides](
		template_model.ToConnectionEntity(derived), template_model.ToConnectionEntity(edited.Clone()))

	var edit ManualConnectionEdit
	edit.Name = edited.Name
	edit.From = edited.From
	edit.To = edited.To
	edit.ConnectionType = edited.ConnectionType
	edit.IsUserAdded = edited.IsUserAdded
	edit.FieldOverrides = overrides
	edit.Unset = unset
	return edit
}

// Apply rebuilds the edited connection on top of derived.
func (this ManualConnectionEdit) Apply(derived template_model.Connection) template_model.Connection {
	connection := template_model.ToConnectionModel(ApplyFieldOverrides(
		template_model.ToConnectionEntity(derived), this.FieldOverrides, this.Unset)).Clone()
	connection.Name = this.Name
	connection.From = this.From
	connection.To = this.To
	connection.ConnectionType = this.ConnectionType
	connection.IsUserAdded = this.IsUserAdded
	return connection
}

// Structure returns a connection carrying only the stored structure fields.
func (this ManualConnectionEdit) Structure() template_model.Connection {
	var connection template_model.Connection
	connection.Name = this.Name
	connection.From = this.From
	connection.To = this.To
	connection.ConnectionType = this.ConnectionType
	connection.IsUserAdded = this.IsUserAdded
	return connection
}

func (this ManualConnectionEdit) Clone() ManualConnectionEdit {
	clone := this
	clone.FieldOverrides = cloneFieldOverrides(this.FieldOverrides)
	clone.Unset = slices.Clone(this.Unset)
	return clone
}
