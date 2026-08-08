package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/config"

// BonusCompositionResultDto is either the composed bonus entries or the
// validation message to show the user; never both.
type BonusCompositionResultDto struct {
	Entries []config.BonusEntry
	Error   string
}
