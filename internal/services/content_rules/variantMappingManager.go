package content_rules

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var (
	resourceObjects            = registry.GetMapObjectResourceValues()
	resourceBankObjects        = registry.GetMapObjectResourceBankValues()
	guardedResourceBankObjects = registry.GetMapObjectT3GuardedResourceBankValues()
)

var (
	// UtopiaVariants enumerates the guard-strength variants of a Dragon Utopia.
	UtopiaVariants = models.NewVariantMapping(
		models.SidMapping{Sid: guardedResourceBankObjects.DragonUtopia, Name: "Dragon Utopia"},
		[]data.Tuple[int, string]{
			data.NewTuple(0, "Small Guard"),
			data.NewTuple(1, "Medium Guard"),
			data.NewTuple(2, "Large Guard"),
			data.NewTuple(3, "Maximum Guard"),
		})

	// PandoraBoxVariants enumerates the reward variants of a Pandora Box.
	PandoraBoxVariants = models.NewVariantMapping(
		models.SidMapping{Sid: resourceObjects.PandoraBox, Name: "Pandora Box"},
		[]data.Tuple[int, string]{
			data.NewTuple(0, "Gold T1 (Low)"),
			data.NewTuple(1, "Gold T2"),
			data.NewTuple(2, "Gold T3"),
			data.NewTuple(3, "Gold T4 (High)"),
			data.NewTuple(4, "Experience T1 (Low)"),
			data.NewTuple(5, "Experience T2"),
			data.NewTuple(6, "Experience T3"),
			data.NewTuple(7, "Experience T4 (High)"),
			data.NewTuple(8, "Units T1 (Low)"),
			data.NewTuple(9, "Units T2"),
			data.NewTuple(10, "Units T3"),
			data.NewTuple(11, "Units T4"),
			data.NewTuple(12, "Units T5"),
			data.NewTuple(13, "Units T6"),
			data.NewTuple(14, "Units T7 (High)"),
			data.NewTuple(15, "All Stats T1 (Low)"),
			data.NewTuple(16, "All Stats T2"),
			data.NewTuple(17, "All Stats T3"),
			data.NewTuple(18, "All Stats T4 (High)"),
			data.NewTuple(19, "Magic School Spells: Daylight"),
			data.NewTuple(20, "Magic School Spells: Nightshade"),
			data.NewTuple(21, "Magic School Spells: Arcane"),
			data.NewTuple(22, "Magic School Spells: Primal"),
			data.NewTuple(23, "Spells T1"),
			data.NewTuple(24, "Spells T2"),
			data.NewTuple(25, "Spells T3"),
			data.NewTuple(26, "Spells T4"),
			data.NewTuple(27, "Spells T5"),
		})

	// MontyHallVariants enumerates the artifact-rarity variants of a Monty Hall.
	MontyHallVariants = models.NewVariantMapping(
		models.SidMapping{Sid: resourceBankObjects.MontyHall, Name: "The Monty Hall"},
		[]data.Tuple[int, string]{
			data.NewTuple(0, "Common Artifact"),
			data.NewTuple(1, "Rare Artifact"),
			data.NewTuple(2, "Epic Artifact"),
			data.NewTuple(3, "Legendary Artifact"),
		})
)

// allVariantMappings preserves declaration order to mirror the C# reflection
// ordering used by GetAllVariantMappings.
var allVariantMappings = []models.VariantMapping{
	UtopiaVariants,
	PandoraBoxVariants,
	MontyHallVariants,
}

// GetVariantsForContent returns one single-entry mapping per variant defined for
// the given content (matched by SID), ordered by variant id. When the content
// has no variants it returns an empty slice.
func GetVariantsForContent(content models.SidMapping) []models.VariantMapping {
	for _, mapping := range allVariantMappings {
		if mapping.Content.Sid != content.Sid {
			continue
		}

		keys := mapping.GetVariantIDsInOrder()
		result := make([]models.VariantMapping, 0, len(keys))
		for _, key := range keys {
			if description, ok := mapping.GetVariantByID(key); ok {
				result = append(result,
					models.NewVariantMapping(content, []data.Tuple[int, string]{data.NewTuple(key, description)}))
			}
		}
		return result
	}

	return []models.VariantMapping{}
}

// GetVariantForContentByID returns the single-entry mapping for the given
// content and variant id, or ok == false when none matches.
func GetVariantForContentByID(content models.SidMapping, variantID int) (models.VariantMapping, bool) {
	for _, mapping := range GetVariantsForContent(content) {
		if _, ok := mapping.GetVariantByID(variantID); ok {
			return mapping, true
		}
	}

	return models.VariantMapping{}, false
}
