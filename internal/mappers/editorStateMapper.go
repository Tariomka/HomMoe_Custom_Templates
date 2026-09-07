package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// EditorStateMapper converts between the persisted editor state and the
// runtime one. The model embeds the entity, so the conversion carries the whole
// group set across in one assignment and only the schema version needs deciding.
type EditorStateMapper struct{}

func NewEditorStateMapper() IEditorStateMapper {
	return &EditorStateMapper{}
}

// NewDefaultEntity is the seed a load decodes over: a key the file omits keeps
// the default instead of collapsing to a zero value. A zero-seeded decode
// cannot tell an absent key from an explicit false or 0, so the seed has to be
// in place before the file is read.
func (this *EditorStateMapper) NewDefaultEntity() editor_state.EditorState {
	return this.ToEntity(editor_state_model.NewDefaultEditorStateModel())
}

func (this *EditorStateMapper) ToEntity(state editor_state_model.EditorState) editor_state.EditorState {
	return editor_state.EditorState{
		TemplateIdentity:    state.TemplateIdentity.TemplateIdentity,
		MapSettings:         state.MapSettings.MapSettings,
		PlayerSettings:      state.PlayerSettings.PlayerSettings,
		NeutralZoneSettings: state.NeutralZoneSettings.NeutralZoneSettings,
		CastleSettings:      state.CastleSettings.CastleSettings,
		GenerationSettings:  state.GenerationSettings.GenerationSettings,
		GameRuleSettings:    state.GameRuleSettings.GameRuleSettings,
		ContentSettings:     editor_state_model.ToContentSettingsEntity(state.ContentSettings),
		ManualEditSettings:  editor_state_model.ToManualEditSettingsEntity(state.ManualEditSettings),
		SchemaVersion:       editor_state.CurrentEditorStateSchemaVersion,
	}
}

func (this *EditorStateMapper) ToModel(entity editor_state.EditorState) editor_state_model.EditorState {
	return editor_state_model.EditorState{
		TemplateIdentity:    editor_state_model.TemplateIdentity{TemplateIdentity: entity.TemplateIdentity},
		MapSettings:         editor_state_model.MapSettings{MapSettings: entity.MapSettings},
		PlayerSettings:      editor_state_model.PlayerSettings{PlayerSettings: entity.PlayerSettings},
		NeutralZoneSettings: editor_state_model.NeutralZoneSettings{NeutralZoneSettings: entity.NeutralZoneSettings},
		CastleSettings:      editor_state_model.CastleSettings{CastleSettings: entity.CastleSettings},
		GenerationSettings:  editor_state_model.GenerationSettings{GenerationSettings: entity.GenerationSettings},
		GameRuleSettings:    editor_state_model.GameRuleSettings{GameRuleSettings: entity.GameRuleSettings},
		ContentSettings:     editor_state_model.ToContentSettingsModel(entity.ContentSettings),
		ManualEditSettings:  editor_state_model.ToManualEditSettingsModel(entity.ManualEditSettings),
		SchemaVersion:       editor_state.CurrentEditorStateSchemaVersion,
	}
}
