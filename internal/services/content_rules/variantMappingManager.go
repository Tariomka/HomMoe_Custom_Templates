package content_rules

import (
	"sort"

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
		models.SidMapping{
			Sid:  guardedResourceBankObjects.DragonUtopia,
			Name: "Dragon Utopia",
		}, map[int]string{
			0: "Small Guard",
			1: "Medium Guard",
			2: "Large Guard",
			3: "Maximum Guard",
		})

	// PandoraBoxVariants enumerates the reward variants of a Pandora Box.
	PandoraBoxVariants = models.NewVariantMapping(
		models.SidMapping{
			Sid:  resourceObjects.PandoraBox,
			Name: "Pandora Box",
		}, map[int]string{
			0:  "Gold T1 (Low)",
			1:  "Gold T2",
			2:  "Gold T3",
			3:  "Gold T4 (High)",
			4:  "Experience T1 (Low)",
			5:  "Experience T2",
			6:  "Experience T3",
			7:  "Experience T4 (High)",
			8:  "Units T1 (Low)",
			9:  "Units T2",
			10: "Units T3",
			11: "Units T4",
			12: "Units T5",
			13: "Units T6",
			14: "Units T7 (High)",
			15: "All Stats T1 (Low)",
			16: "All Stats T2",
			17: "All Stats T3",
			18: "All Stats T4 (High)",
			19: "Magic School Spells: Daylight",
			20: "Magic School Spells: Nightshade",
			21: "Magic School Spells: Arcane",
			22: "Magic School Spells: Primal",
			23: "Spells T1",
			24: "Spells T2",
			25: "Spells T3",
			26: "Spells T4",
			27: "Spells T5",
		})

	// MontyHallVariants enumerates the artifact-rarity variants of a Monty Hall.
	MontyHallVariants = models.NewVariantMapping(
		models.SidMapping{
			Sid:  resourceBankObjects.MontyHall,
			Name: "The Monty Hall",
		}, map[int]string{
			0: "Common Artifact",
			1: "Rare Artifact",
			2: "Epic Artifact",
			3: "Legendary Artifact",
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
		keys := make([]int, 0, len(mapping.Variants))
		for key := range mapping.Variants {
			keys = append(keys, key)
		}
		sort.Ints(keys)

		result := make([]models.VariantMapping, 0, len(keys))
		for _, key := range keys {
			result = append(result, models.NewVariantMapping(content, map[int]string{
				key: mapping.Variants[key],
			}))
		}
		return result
	}
	return []models.VariantMapping{}
}

// GetVariantForContentByID returns the single-entry mapping for the given
// content and variant id, or ok == false when none matches.
func GetVariantForContentByID(content models.SidMapping, variantID int) (models.VariantMapping, bool) {
	for _, mapping := range GetVariantsForContent(content) {
		if _, ok := mapping.Variants[variantID]; ok {
			return mapping, true
		}
	}
	return models.VariantMapping{}, false
}
