package editor_state

import "encoding/json"

// CurrentEditorStateSchemaVersion is stamped into every editor state written
// today. Files written before the version existed carry no key at all, decode
// as 0, and are migrated on load.
const CurrentEditorStateSchemaVersion = 1

// EditorStateEntity is the .gen.json file on disk: the nine behaviour-free
// field groups plus the schema version. The groups are embedded anonymously so
// encoding/json keeps emitting a single flat object with schemaVersion as a
// sibling of the settings keys.
//
//nolint:recvcheck // a JSON codec pair is inherently mixed: MarshalJSON must take a value so both a state and a *state serialise through it, UnmarshalJSON must take a pointer to merge into the receiver.
type EditorStateEntity struct {
	TemplateIdentity
	MapSettings
	PlayerSettings
	NeutralZoneSettings
	CastleSettings
	GenerationSettings
	GameRuleSettings
	ContentSettings
	ManualEditSettings

	SchemaVersion int `json:"schemaVersion"`
}

// editorStateFields has EditorStateEntity's layout but none of its methods, so
// the codecs below cannot re-enter themselves through it.
type editorStateFields EditorStateEntity

// MarshalJSON always stamps the current schema version, whatever the receiver
// happens to carry, so a state loaded from a legacy file is written back at the
// version the writer actually produces.
func (this EditorStateEntity) MarshalJSON() ([]byte, error) {
	this.SchemaVersion = CurrentEditorStateSchemaVersion

	return json.Marshal(editorStateFields(this))
}

// UnmarshalJSON overlays the file onto the receiver rather than replacing it, so
// a key the file omits keeps whatever the caller seeded. That is what lets an
// absent setting fall back to its default instead of collapsing to a zero value.
func (this *EditorStateEntity) UnmarshalJSON(data []byte) error {
	fields := editorStateFields(*this)
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	*this = EditorStateEntity(fields)
	this.migrateSchemaVersion()

	return nil
}

func (this *EditorStateEntity) migrateSchemaVersion() {
	if this.SchemaVersion == 0 {
		this.SchemaVersion = CurrentEditorStateSchemaVersion
	}
}
