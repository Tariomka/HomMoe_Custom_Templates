package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// EditorStateEntityMapper converts between the persisted editor state and the
// runtime one. Both carry the same nine field groups, so the conversion is a
// regrouping rather than a copy: the slices inside are shared with the argument.
type EditorStateEntityMapper struct{}

func NewEditorStateEntityMapper() IEditorStateEntityMapper {
	return &EditorStateEntityMapper{}
}

func (this *EditorStateEntityMapper) ToEntity(state editor_state_model.EditorState) editor_state.EditorStateEntity {
	return editor_state.EditorStateEntity{
		SchemaVersion:       editor_state.CurrentEditorStateSchemaVersion,
		TemplateIdentity:    state.TemplateIdentity,
		MapSettings:         state.MapSettings,
		PlayerSettings:      state.PlayerSettings,
		NeutralZoneSettings: state.NeutralZoneSettings,
		CastleSettings:      state.CastleSettings,
		GenerationSettings:  state.GenerationSettings,
		GameRuleSettings:    state.GameRuleSettings,
		ContentSettings:     state.ContentSettings,
		ManualEditSettings:  state.ManualEditSettings,
	}
}

func (this *EditorStateEntityMapper) ToModel(entity editor_state.EditorStateEntity) editor_state_model.EditorState {
	return editor_state_model.EditorState{
		TemplateIdentity:    entity.TemplateIdentity,
		MapSettings:         entity.MapSettings,
		PlayerSettings:      entity.PlayerSettings,
		NeutralZoneSettings: entity.NeutralZoneSettings,
		CastleSettings:      entity.CastleSettings,
		GenerationSettings:  entity.GenerationSettings,
		GameRuleSettings:    entity.GameRuleSettings,
		ContentSettings:     entity.ContentSettings,
		ManualEditSettings:  entity.ManualEditSettings,
	}
}
