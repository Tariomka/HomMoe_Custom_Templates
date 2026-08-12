package bonuses

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// IBonusEntryService composes and validates game-start bonus entries for the
// bonus picker.
type IBonusEntryService interface {
	// DescribeExistingBonuses summarises an already-composed bonus list for
	// duplicate detection and spell pre-exclusion.
	DescribeExistingBonuses(existing []config.BonusEntry) dtos.ExistingBonusesDto

	// BuildBonusEntries turns filled-in composer form values into bonus
	// entries, or reports why the form is not usable yet.
	BuildBonusEntries(request dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto

	// FilterNewBonusEntries drops entries whose hash is already present.
	FilterNewBonusEntries(entries []config.BonusEntry, existingKeys map[string]bool) []config.BonusEntry

	// GetSpellCountLabel renders the picked-spell summary line.
	GetSpellCountLabel(count int) string
}
