package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_generator"
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

func (this *EditorStateMapper) NewDefaultEntity() editor_state.EditorState {
	return this.ToEntity(editor_state_model.NewDefaultEditorStateModel())
}

func (this *EditorStateMapper) ToEntity(state editor_state_model.EditorState) editor_state.EditorState {
	generatorVersion := common_generator.GeneratorOutputVersion
	if state.NeedsDerivation() {
		generatorVersion = 0
	}

	return editor_state.EditorState{
		SchemaVersion:    editor_state.CurrentEditorStateSchemaVersion,
		GeneratorVersion: generatorVersion,
		TemplateIdentity: state.TemplateIdentity.TemplateIdentity,
		MapSettings:      state.MapSettings.MapSettings,
		Players:          state.PlayerSettings.PlayerSettings,
		NeutralZones:     state.NeutralZoneSettings.NeutralZoneSettings,
		Castles:          state.CastleSettings.CastleSettings,
		Generation:       state.GenerationSettings.GenerationSettings,
		GameRules:        state.GameRuleSettings.GameRuleSettings,
		Content:          editor_state_model.ToContentSettingsEntity(state.ContentSettings),
		ManualEdits:      editor_state_model.ToManualEditSettingsEntity(state.ManualEditSettings),
	}
}

func (this *EditorStateMapper) ToModel(entity editor_state.EditorState) editor_state_model.EditorState {
	return editor_state_model.EditorState{
		TemplateIdentity:    editor_state_model.TemplateIdentity{TemplateIdentity: entity.TemplateIdentity},
		MapSettings:         editor_state_model.MapSettings{MapSettings: entity.MapSettings},
		PlayerSettings:      editor_state_model.PlayerSettings{PlayerSettings: entity.Players},
		NeutralZoneSettings: editor_state_model.NeutralZoneSettings{NeutralZoneSettings: entity.NeutralZones},
		CastleSettings:      editor_state_model.CastleSettings{CastleSettings: entity.Castles},
		GenerationSettings:  editor_state_model.GenerationSettings{GenerationSettings: entity.Generation},
		GameRuleSettings:    editor_state_model.GameRuleSettings{GameRuleSettings: entity.GameRules},
		ContentSettings:     editor_state_model.ToContentSettingsModel(entity.Content),
		ManualEditSettings:  editor_state_model.ToManualEditSettingsModel(entity.ManualEdits),
		SchemaVersion:       editor_state.CurrentEditorStateSchemaVersion,
		GeneratorVersion:    entity.GeneratorVersion,
	}
}
