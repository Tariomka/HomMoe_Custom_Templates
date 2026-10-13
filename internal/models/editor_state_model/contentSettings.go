package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentSettings struct {
	BannedItems    []BanEntry
	BannedMagics   []BanEntry
	ValueOverrides []ValueOverrideEntry
	Bonuses        []BonusEntry

	PlayerZoneContentRows    []ZoneContentRow
	LowestNeutralContentRows []ZoneContentRow
	LowNeutralContentRows    []ZoneContentRow
	MediumNeutralContentRows []ZoneContentRow
	HighNeutralContentRows   []ZoneContentRow
	HubZoneContentRows       []ZoneContentRow
}

func ToContentSettingsModel(entity editor_state.ContentSettings) ContentSettings {
	playerRows := ToZoneContentRowModels(entity.PlayerZoneContentRows)
	if entity.PlayerZoneContentRows == nil {
		var known bool
		if playerRows, known = GetDefaultPlayerZoneContentRowsForVersion(entity.ContentDefaultsVersion); !known {
			playerRows = GetDefaultPlayerZoneContentRows()
		}
	}

	return ContentSettings{
		BannedItems:    ToBanEntryModels(entity.BannedItems),
		BannedMagics:   ToBanEntryModels(entity.BannedMagics),
		ValueOverrides: ToValueOverrideEntryModels(entity.ValueOverrides),
		Bonuses:        ToBonusEntryModels(entity.Bonuses),

		PlayerZoneContentRows:    playerRows,
		LowestNeutralContentRows: ToZoneContentRowModels(entity.LowestNeutralContentRows),
		LowNeutralContentRows:    ToZoneContentRowModels(entity.LowNeutralContentRows),
		MediumNeutralContentRows: ToZoneContentRowModels(entity.MediumNeutralContentRows),
		HighNeutralContentRows:   ToZoneContentRowModels(entity.HighNeutralContentRows),
		HubZoneContentRows:       ToZoneContentRowModels(entity.HubZoneContentRows),
	}
}

func ToContentSettingsEntity(model ContentSettings) editor_state.ContentSettings {
	playerRows := ToZoneContentRowEntities(model.PlayerZoneContentRows)
	switch {
	case helpers.EqualFunc(model.PlayerZoneContentRows, GetDefaultPlayerZoneContentRows(), contentRowsEqual):
		playerRows = nil
	case len(playerRows) == 0:
		playerRows = []editor_state.ZoneContentRow{}
	}

	return editor_state.ContentSettings{
		ContentDefaultsVersion: editor_state.CurrentContentDefaultsVersion,
		BannedItems:            ToBanEntryEntities(model.BannedItems),
		BannedMagics:           ToBanEntryEntities(model.BannedMagics),
		ValueOverrides:         ToValueOverrideEntryEntities(model.ValueOverrides),
		Bonuses:                ToBonusEntryEntities(model.Bonuses),

		PlayerZoneContentRows:    playerRows,
		LowestNeutralContentRows: ToZoneContentRowEntities(model.LowestNeutralContentRows),
		LowNeutralContentRows:    ToZoneContentRowEntities(model.LowNeutralContentRows),
		MediumNeutralContentRows: ToZoneContentRowEntities(model.MediumNeutralContentRows),
		HighNeutralContentRows:   ToZoneContentRowEntities(model.HighNeutralContentRows),
		HubZoneContentRows:       ToZoneContentRowEntities(model.HubZoneContentRows),
	}
}
