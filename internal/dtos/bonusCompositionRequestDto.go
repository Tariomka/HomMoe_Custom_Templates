package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/config"

// BonusCompositionRequestDto carries the raw bonus-composer form values. The
// text fields are unparsed on purpose - validating them is the service's job.
type BonusCompositionRequestDto struct {
	PresetType     config.BonusPresetType
	ReceiverFilter string
	SelectedSpells []string
	MakeSpellsFree bool
	MultiplierText string
	MovementText   string
	ItemText       string
	ResourceText   string
}
