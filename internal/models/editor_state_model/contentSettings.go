package editor_state_model

import "github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"

type ContentSettings struct {
	BannedItems        string
	BannedMagics       string
	ValueOverridesText string
	Bonuses            []BonusEntry

	PlayerZoneContentRows    []ZoneContentRow
	LowestNeutralContentRows []ZoneContentRow
	LowNeutralContentRows    []ZoneContentRow
	MediumNeutralContentRows []ZoneContentRow
	HighNeutralContentRows   []ZoneContentRow
	HubZoneContentRows       []ZoneContentRow
}

func ToContentSettingsModel(entity editor_state.ContentSettings) ContentSettings {
	return ContentSettings{
		BannedItems:        entity.BannedItems,
		BannedMagics:       entity.BannedMagics,
		ValueOverridesText: entity.ValueOverridesText,
		Bonuses:            ToBonusEntryModels(entity.Bonuses),

		PlayerZoneContentRows:    ToZoneContentRowModels(entity.PlayerZoneContentRows),
		LowestNeutralContentRows: ToZoneContentRowModels(entity.LowestNeutralContentRows),
		LowNeutralContentRows:    ToZoneContentRowModels(entity.LowNeutralContentRows),
		MediumNeutralContentRows: ToZoneContentRowModels(entity.MediumNeutralContentRows),
		HighNeutralContentRows:   ToZoneContentRowModels(entity.HighNeutralContentRows),
		HubZoneContentRows:       ToZoneContentRowModels(entity.HubZoneContentRows),
	}
}

func ToContentSettingsEntity(model ContentSettings) editor_state.ContentSettings {
	return editor_state.ContentSettings{
		BannedItems:        model.BannedItems,
		BannedMagics:       model.BannedMagics,
		ValueOverridesText: model.ValueOverridesText,
		Bonuses:            ToBonusEntryEntities(model.Bonuses),

		PlayerZoneContentRows:    ToZoneContentRowEntities(model.PlayerZoneContentRows),
		LowestNeutralContentRows: ToZoneContentRowEntities(model.LowestNeutralContentRows),
		LowNeutralContentRows:    ToZoneContentRowEntities(model.LowNeutralContentRows),
		MediumNeutralContentRows: ToZoneContentRowEntities(model.MediumNeutralContentRows),
		HighNeutralContentRows:   ToZoneContentRowEntities(model.HighNeutralContentRows),
		HubZoneContentRows:       ToZoneContentRowEntities(model.HubZoneContentRows),
	}
}
