package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// EditorStateMapper converts the editor state across the app/ <-> internal/
// crossing.
//
// Every conversion below is a plain struct conversion rather than a field list.
// The DTO groups mirror the entity groups exactly and Go ignores struct tags
// when deciding conversion identity, so adding a field to one side without the
// other is a compile error - which is the whole reason the DTO is grouped
// rather than flat.
type EditorStateMapper struct{}

func NewEditorStateMapper() IEditorStateMapper {
	return &EditorStateMapper{}
}

func (this *EditorStateMapper) ToDto(
	state editor_state_model.EditorState) editor_state_dto.EditorStateDto {
	return editor_state_dto.EditorStateDto{
		TemplateIdentityDto:    editor_state_dto.TemplateIdentityDto(state.TemplateIdentity),
		MapSettingsDto:         editor_state_dto.MapSettingsDto(state.MapSettings),
		PlayerSettingsDto:      editor_state_dto.PlayerSettingsDto(state.PlayerSettings),
		NeutralZoneSettingsDto: editor_state_dto.NeutralZoneSettingsDto(state.NeutralZoneSettings),
		CastleSettingsDto:      editor_state_dto.CastleSettingsDto(state.CastleSettings),
		GenerationSettingsDto:  editor_state_dto.GenerationSettingsDto(state.GenerationSettings),
		GameRuleSettingsDto:    editor_state_dto.GameRuleSettingsDto(state.GameRuleSettings),
		ContentSettingsDto:     editor_state_dto.ContentSettingsDto(state.ContentSettings),
		ManualEditSettingsDto:  editor_state_dto.ManualEditSettingsDto(state.ManualEditSettings),
	}
}

// ToModel stamps the current schema version because the DTO deliberately does
// not carry one: a state that came in over the crossing is by definition the
// shape this build understands.
func (this *EditorStateMapper) ToModel(
	dto editor_state_dto.EditorStateDto) editor_state_model.EditorState {
	return editor_state_model.EditorState{
		TemplateIdentity:    editor_state.TemplateIdentity(dto.TemplateIdentityDto),
		MapSettings:         editor_state.MapSettings(dto.MapSettingsDto),
		PlayerSettings:      editor_state.PlayerSettings(dto.PlayerSettingsDto),
		NeutralZoneSettings: editor_state.NeutralZoneSettings(dto.NeutralZoneSettingsDto),
		CastleSettings:      editor_state.CastleSettings(dto.CastleSettingsDto),
		GenerationSettings:  editor_state.GenerationSettings(dto.GenerationSettingsDto),
		GameRuleSettings:    editor_state.GameRuleSettings(dto.GameRuleSettingsDto),
		ContentSettings:     editor_state.ContentSettings(dto.ContentSettingsDto),
		ManualEditSettings:  editor_state.ManualEditSettings(dto.ManualEditSettingsDto),
		SchemaVersion:       editor_state.CurrentEditorStateSchemaVersion}
}

// ToDtoPointer keeps nil meaning "absent": the regeneration request uses a nil
// Previous or Next to distinguish a first generation from a steady state.
func (this *EditorStateMapper) ToDtoPointer(
	state *editor_state_model.EditorState) *editor_state_dto.EditorStateDto {
	if state == nil {
		return nil
	}

	return new(this.ToDto(*state))
}

func (this *EditorStateMapper) ToModelPointer(
	dto *editor_state_dto.EditorStateDto) *editor_state_model.EditorState {
	if dto == nil {
		return nil
	}

	return new(this.ToModel(*dto))
}

func (this *EditorStateMapper) ToCastleSettingChangesDto(
	changes editor_state_model.CastleSettingChanges) editor_state_dto.CastleSettingChangesDto {
	return editor_state_dto.CastleSettingChangesDto(changes)
}

func (this *EditorStateMapper) ToCastleSettingChangesModel(
	dto editor_state_dto.CastleSettingChangesDto) editor_state_model.CastleSettingChanges {
	return editor_state_model.CastleSettingChanges(dto)
}
