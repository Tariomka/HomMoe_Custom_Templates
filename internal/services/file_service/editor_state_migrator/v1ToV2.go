package editor_state_migrator

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

// MigrateToV2 upgrades a decoded v1 (or v0) editor state to the current shape.
// Every settings group crosses by plain struct conversion, and that is the point
// of the frozen snapshot: adding or renaming a field on a current group breaks
// this file at compile time instead of silently reinterpreting old files.
func MigrateToV2(legacy editor_state_v1.EditorState) editor_state.EditorState {
	return editor_state.EditorState{
		TemplateIdentity:    editor_state.TemplateIdentity(legacy.TemplateIdentity),
		MapSettings:         editor_state.MapSettings(legacy.MapSettings),
		PlayerSettings:      editor_state.PlayerSettings(legacy.PlayerSettings),
		NeutralZoneSettings: editor_state.NeutralZoneSettings(legacy.NeutralZoneSettings),
		CastleSettings:      editor_state.CastleSettings(legacy.CastleSettings),
		GenerationSettings:  editor_state.GenerationSettings(legacy.GenerationSettings),
		GameRuleSettings:    editor_state.GameRuleSettings(legacy.GameRuleSettings),
		ContentSettings:     toContentSettings(legacy.ContentSettings),
		ManualEditSettings:  toManualEditSettings(legacy.ManualEditSettings),
		SchemaVersion:       editor_state.CurrentEditorStateSchemaVersion,
	}
}

// newV1Seed lowers the seeded defaults into the v1 shape so a v1 decode keeps a
// default for every key the file omits. Manual edits are deliberately not
// carried: defaults never hold any, and a v1 file's manual edits replace them
// wholesale.
func newV1Seed(seed editor_state.EditorState) editor_state_v1.EditorState {
	return editor_state_v1.EditorState{
		TemplateIdentity:    editor_state_v1.TemplateIdentity(seed.TemplateIdentity),
		MapSettings:         editor_state_v1.MapSettings(seed.MapSettings),
		PlayerSettings:      editor_state_v1.PlayerSettings(seed.PlayerSettings),
		NeutralZoneSettings: editor_state_v1.NeutralZoneSettings(seed.NeutralZoneSettings),
		CastleSettings:      editor_state_v1.CastleSettings(seed.CastleSettings),
		GenerationSettings:  editor_state_v1.GenerationSettings(seed.GenerationSettings),
		GameRuleSettings:    editor_state_v1.GameRuleSettings(seed.GameRuleSettings),
		ContentSettings:     toV1ContentSettings(seed.ContentSettings),
	}
}

func toContentSettings(legacy editor_state_v1.ContentSettings) editor_state.ContentSettings {
	return editor_state.ContentSettings{
		BannedItems:              legacy.BannedItems,
		BannedMagics:             legacy.BannedMagics,
		ValueOverridesText:       legacy.ValueOverridesText,
		Bonuses:                  helpers.MapSlice(legacy.Bonuses, toBonusEntry),
		PlayerZoneContentRows:    helpers.MapSlice(legacy.PlayerZoneContentRows, toZoneContentRow),
		LowestNeutralContentRows: helpers.MapSlice(legacy.LowestNeutralContentRows, toZoneContentRow),
		LowNeutralContentRows:    helpers.MapSlice(legacy.LowNeutralContentRows, toZoneContentRow),
		MediumNeutralContentRows: helpers.MapSlice(legacy.MediumNeutralContentRows, toZoneContentRow),
		HighNeutralContentRows:   helpers.MapSlice(legacy.HighNeutralContentRows, toZoneContentRow),
		HubZoneContentRows:       helpers.MapSlice(legacy.HubZoneContentRows, toZoneContentRow),
	}
}

func toV1ContentSettings(seed editor_state.ContentSettings) editor_state_v1.ContentSettings {
	return editor_state_v1.ContentSettings{
		BannedItems:              seed.BannedItems,
		BannedMagics:             seed.BannedMagics,
		ValueOverridesText:       seed.ValueOverridesText,
		Bonuses:                  helpers.MapSlice(seed.Bonuses, toV1BonusEntry),
		PlayerZoneContentRows:    helpers.MapSlice(seed.PlayerZoneContentRows, toV1ZoneContentRow),
		LowestNeutralContentRows: helpers.MapSlice(seed.LowestNeutralContentRows, toV1ZoneContentRow),
		LowNeutralContentRows:    helpers.MapSlice(seed.LowNeutralContentRows, toV1ZoneContentRow),
		MediumNeutralContentRows: helpers.MapSlice(seed.MediumNeutralContentRows, toV1ZoneContentRow),
		HighNeutralContentRows:   helpers.MapSlice(seed.HighNeutralContentRows, toV1ZoneContentRow),
		HubZoneContentRows:       helpers.MapSlice(seed.HubZoneContentRows, toV1ZoneContentRow),
	}
}

func toBonusEntry(legacy editor_state_v1.BonusEntry) editor_state.BonusEntry {
	return editor_state.BonusEntry{
		PresetType:     editor_state.BonusPresetType(legacy.PresetType),
		ReceiverFilter: legacy.ReceiverFilter,
		Param:          legacy.Param,
		Param2:         legacy.Param2,
	}
}

func toV1BonusEntry(seed editor_state.BonusEntry) editor_state_v1.BonusEntry {
	return editor_state_v1.BonusEntry{
		PresetType:     editor_state_v1.BonusPresetType(seed.PresetType),
		ReceiverFilter: seed.ReceiverFilter,
		Param:          seed.Param,
		Param2:         seed.Param2,
	}
}

func toZoneContentRow(legacy editor_state_v1.ZoneContentRow) editor_state.ZoneContentRow {
	return editor_state.ZoneContentRow{
		Sid:     legacy.Sid,
		Count:   legacy.Count,
		IsGroup: legacy.IsGroup,
		IsMine:  legacy.IsMine,
		Rules:   helpers.MapSlice(legacy.Rules, toContentRuleRow),
	}
}

func toV1ZoneContentRow(seed editor_state.ZoneContentRow) editor_state_v1.ZoneContentRow {
	return editor_state_v1.ZoneContentRow{
		Sid:     seed.Sid,
		Count:   seed.Count,
		IsGroup: seed.IsGroup,
		IsMine:  seed.IsMine,
		Rules:   helpers.MapSlice(seed.Rules, toV1ContentRuleRow),
	}
}

func toContentRuleRow(legacy editor_state_v1.ContentRuleRow) editor_state.ContentRuleRow {
	return editor_state.ContentRuleRow(legacy)
}

func toV1ContentRuleRow(seed editor_state.ContentRuleRow) editor_state_v1.ContentRuleRow {
	return editor_state_v1.ContentRuleRow(seed)
}

func toManualEditSettings(legacy editor_state_v1.ManualEditSettings) editor_state.ManualEditSettings {
	return editor_state.ManualEditSettings{
		ManualZones:       helpers.MapSlice(legacy.ManualZones, toManualZoneSave),
		ManualConnections: helpers.MapSlice(legacy.ManualConnections, toManualConnectionSave),
	}
}

// toManualZoneSave leaves GeneratorPosition and GeneratorRing nil: v1 never
// persisted them, and nil correctly reads as "the stamp was never taken".
func toManualZoneSave(legacy editor_state_v1.ManualZoneSave) editor_state.ManualZoneSave {
	return editor_state.ManualZoneSave{
		Zone:           legacy.Zone,
		ManualPosition: toVector(legacy.ManualPosition),
		Quality:        legacy.Quality,
	}
}

func toManualConnectionSave(
	legacy editor_state_v1.ManualConnectionSave) editor_state.ManualConnectionSave {
	return editor_state.ManualConnectionSave(legacy)
}

func toVector(position *[2]float64) *data.Vec2[float64] {
	if position == nil {
		return nil
	}

	return new(data.NewVec2(position[0], position[1]))
}
