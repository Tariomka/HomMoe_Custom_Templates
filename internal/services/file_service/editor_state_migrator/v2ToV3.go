package editor_state_migrator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v2"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// legacyContentDefaultsVersion is the default-rows table every pre-v3 file was
// written against.
const legacyContentDefaultsVersion = 1

// MigrateToV3 upgrades a decoded v2 editor state to the current sectioned
// shape. The ban and override text is parsed here, once: blank lines are
// dropped, order and duplicates are kept, and each malformed override line is
// dropped with a warning. The manual snapshot is kept whole, in the
// all-overridden form, with generator version 0 so it is converted at the
// first regeneration.
func MigrateToV3(legacy editor_state_v2.EditorState) (editor_state.EditorState, []string) {
	valueOverrides, warnings := parseValueOverrides(legacy.ValueOverridesText)
	return editor_state.EditorState{
		SchemaVersion:    editor_state.CurrentEditorStateSchemaVersion,
		GeneratorVersion: 0,
		TemplateIdentity: editor_state.TemplateIdentity(legacy.TemplateIdentity),
		MapSettings:      editor_state.MapSettings(legacy.MapSettings),
		Players:          editor_state.PlayerSettings(legacy.PlayerSettings),
		NeutralZones:     editor_state.NeutralZoneSettings(legacy.NeutralZoneSettings),
		Castles:          editor_state.CastleSettings(legacy.CastleSettings),
		Generation:       editor_state.GenerationSettings(legacy.GenerationSettings),
		GameRules:        editor_state.GameRuleSettings(legacy.GameRuleSettings),
		Content: editor_state.ContentSettings{
			ContentDefaultsVersion:   legacyContentDefaultsVersion,
			BannedItems:              parseBanEntries(legacy.BannedItems),
			BannedMagics:             parseBanEntries(legacy.BannedMagics),
			ValueOverrides:           valueOverrides,
			Bonuses:                  helpers.MapSlice(legacy.Bonuses, fromV2BonusEntry),
			PlayerZoneContentRows:    helpers.MapSlice(legacy.PlayerZoneContentRows, fromV2ZoneContentRow),
			LowestNeutralContentRows: helpers.MapSlice(legacy.LowestNeutralContentRows, fromV2ZoneContentRow),
			LowNeutralContentRows:    helpers.MapSlice(legacy.LowNeutralContentRows, fromV2ZoneContentRow),
			MediumNeutralContentRows: helpers.MapSlice(legacy.MediumNeutralContentRows, fromV2ZoneContentRow),
			HighNeutralContentRows:   helpers.MapSlice(legacy.HighNeutralContentRows, fromV2ZoneContentRow),
			HubZoneContentRows:       helpers.MapSlice(legacy.HubZoneContentRows, fromV2ZoneContentRow),
		},
		ManualEdits: fromV2ManualEditSettings(legacy.ManualEditSettings),
	}, warnings
}

// newV2Seed lowers the seeded defaults into the v2 shape so a v2 decode keeps a
// default for every key the file omits.
func newV2Seed(seed editor_state.EditorState) editor_state_v2.EditorState {
	return editor_state_v2.EditorState{
		TemplateIdentity:         editor_state_v2.TemplateIdentity(seed.TemplateIdentity),
		MapSettings:              editor_state_v2.MapSettings(seed.MapSettings),
		PlayerSettings:           editor_state_v2.PlayerSettings(seed.Players),
		NeutralZoneSettings:      editor_state_v2.NeutralZoneSettings(seed.NeutralZones),
		CastleSettings:           editor_state_v2.CastleSettings(seed.Castles),
		GenerationSettings:       editor_state_v2.GenerationSettings(seed.Generation),
		GameRuleSettings:         editor_state_v2.GameRuleSettings(seed.GameRules),
		BannedItems:              formatBanEntries(seed.Content.BannedItems),
		BannedMagics:             formatBanEntries(seed.Content.BannedMagics),
		ValueOverridesText:       formatValueOverrides(seed.Content.ValueOverrides),
		Bonuses:                  helpers.MapSlice(seed.Content.Bonuses, toV2BonusEntry),
		PlayerZoneContentRows:    helpers.MapSlice(seed.Content.PlayerZoneContentRows, toV2ZoneContentRow),
		LowestNeutralContentRows: helpers.MapSlice(seed.Content.LowestNeutralContentRows, toV2ZoneContentRow),
		LowNeutralContentRows:    helpers.MapSlice(seed.Content.LowNeutralContentRows, toV2ZoneContentRow),
		MediumNeutralContentRows: helpers.MapSlice(seed.Content.MediumNeutralContentRows, toV2ZoneContentRow),
		HighNeutralContentRows:   helpers.MapSlice(seed.Content.HighNeutralContentRows, toV2ZoneContentRow),
		HubZoneContentRows:       helpers.MapSlice(seed.Content.HubZoneContentRows, toV2ZoneContentRow),
	}
}

func parseBanEntries(text string) []editor_state.BanEntry {
	var entries []editor_state.BanEntry
	for line := range strings.SplitSeq(text, "\n") {
		if sid := strings.TrimSpace(line); sid != "" {
			entries = append(entries, editor_state.BanEntry{Sid: sid})
		}
	}
	return entries
}

// parseValueOverrides keeps the generator's former "sid=guardValue" rules: the
// SID ends at the first '=', both halves are trimmed, the value must be an
// integer. Line numbers count every line, blank ones included.
func parseValueOverrides(text string) ([]editor_state.ValueOverrideEntry, []string) {
	var entries []editor_state.ValueOverrideEntry
	var warnings []string
	for index, rawLine := range strings.Split(text, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		equals := strings.Index(line, "=")
		if equals <= 0 {
			warnings = append(warnings,
				fmt.Sprintf("value override line %d '%s' is not sid=value and was dropped", index+1, line))
			continue
		}

		guardValue, err := strconv.Atoi(strings.TrimSpace(line[equals+1:]))
		if err != nil {
			warnings = append(warnings,
				fmt.Sprintf("value override line %d '%s' has a non-numeric value and was dropped", index+1, line))
			continue
		}

		entries = append(entries, editor_state.ValueOverrideEntry{
			Sid:        strings.TrimSpace(line[:equals]),
			GuardValue: guardValue,
		})
	}
	return entries, warnings
}

func formatBanEntries(entries []editor_state.BanEntry) string {
	sids := make([]string, 0, len(entries))
	for _, entry := range entries {
		sids = append(sids, entry.Sid)
	}
	return strings.Join(sids, "\n")
}

func formatValueOverrides(entries []editor_state.ValueOverrideEntry) string {
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, fmt.Sprintf("%s=%d", entry.Sid, entry.GuardValue))
	}
	return strings.Join(lines, "\n")
}

func fromV2BonusEntry(legacy editor_state_v2.BonusEntry) editor_state.BonusEntry {
	return editor_state.BonusEntry{
		PresetType:     editor_state.BonusPresetType(legacy.PresetType),
		ReceiverFilter: legacy.ReceiverFilter,
		Param:          legacy.Param,
		Param2:         legacy.Param2,
	}
}

func toV2BonusEntry(seed editor_state.BonusEntry) editor_state_v2.BonusEntry {
	return editor_state_v2.BonusEntry{
		PresetType:     editor_state_v2.BonusPresetType(seed.PresetType),
		ReceiverFilter: seed.ReceiverFilter,
		Param:          seed.Param,
		Param2:         seed.Param2,
	}
}

func fromV2ZoneContentRow(legacy editor_state_v2.ZoneContentRow) editor_state.ZoneContentRow {
	return editor_state.ZoneContentRow{
		Sid:     legacy.Sid,
		Count:   legacy.Count,
		IsGroup: legacy.IsGroup,
		IsMine:  legacy.IsMine,
		Rules: helpers.MapSlice(legacy.Rules, func(rule editor_state_v2.ContentRuleRow) editor_state.ContentRuleRow {
			return editor_state.ContentRuleRow(rule)
		}),
	}
}

func toV2ZoneContentRow(seed editor_state.ZoneContentRow) editor_state_v2.ZoneContentRow {
	return editor_state_v2.ZoneContentRow{
		Sid:     seed.Sid,
		Count:   seed.Count,
		IsGroup: seed.IsGroup,
		IsMine:  seed.IsMine,
		Rules: helpers.MapSlice(seed.Rules, func(rule editor_state.ContentRuleRow) editor_state_v2.ContentRuleRow {
			return editor_state_v2.ContentRuleRow(rule)
		}),
	}
}

// fromV2ManualEditSettings rebuilds the v2 snapshot as live zones and
// connections and writes them in the all-overridden form, so the v3 encoding
// has a single implementation.
func fromV2ManualEditSettings(legacy editor_state_v2.ManualEditSettings) editor_state.ManualEditSettings {
	zones := helpers.MapSlice(legacy.ManualZones, func(save editor_state_v2.ManualZoneSave) template_model.Zone {
		zone := template_model.ToZoneModel(save.Zone)
		zone.GeneratorPosition = save.GeneratorPosition
		zone.GeneratorRing = save.GeneratorRing
		zone.ManualPosition = save.ManualPosition
		zone.Quality = helpers.MapPointer(save.Quality, func(ordinal int8) neutral_zone.Quality {
			return neutral_zone.Quality(ordinal)
		})
		return zone
	})
	connections := helpers.MapSlice(legacy.ManualConnections,
		func(save editor_state_v2.ManualConnectionSave) template_model.Connection {
			connection := template_model.ToConnectionModel(save.Connection)
			connection.IsUserAdded = save.IsUserAdded
			return connection
		})
	return editor_state_model.ToManualEditSettingsEntity(
		editor_state_model.ManualEditSettingsFromSnapshot(zones, connections))
}
