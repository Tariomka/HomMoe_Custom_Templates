package content_rules

import (
	"fmt"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// Rule metadata constants for the variant rule.
const (
	RuleVariantName        = "Variant"
	RuleVariantDescription = "Forces the content item to spawn a specific variant."
	RuleVariantMarker      = ""
)

// RuleVariant forces the content item to spawn a specific variant.
type RuleVariant struct {
	Mapping   models.VariantMapping
	VariantId int
}

// NewRuleVariant creates a variant rule. When no mapping is supplied it falls
// back to the dragon-utopia mapping (matching the C# dummy default). When no
// variant id is supplied it uses the lowest defined variant id for determinism.
// It returns an error when the resolved id is not present in the mapping.
func NewRuleVariant(mapping *models.VariantMapping, variantId *int) (*RuleVariant, error) {
	resolved := UtopiaVariants
	if mapping != nil {
		resolved = *mapping
	}

	var id int
	if variantId != nil {
		id = *variantId
	} else {
		id = smallestVariantKey(resolved.Variants)
	}

	if _, ok := resolved.Variants[id]; !ok {
		return nil, fmt.Errorf("selected variant id %d is not present in the provided variant mapping", id)
	}
	return &RuleVariant{Mapping: resolved, VariantId: id}, nil
}

func (this *RuleVariant) Name() string        { return RuleVariantName }
func (this *RuleVariant) Description() string { return RuleVariantDescription }
func (this *RuleVariant) Marker() string      { return RuleVariantMarker }

func (this *RuleVariant) DisplayText() string {
	if description, ok := this.Mapping.Variants[this.VariantId]; ok {
		return fmt.Sprintf("%s: %s", this.Name(), description)
	}
	return fmt.Sprintf("%s: Unforeseen Error", this.Name())
}

func (this *RuleVariant) Apply(item *entities.MandatoryContentItem) {
	id := this.VariantId
	item.Variant = &id
}

func (this *RuleVariant) SerializeToRowSave() models.ContentRuleRowSave {
	id := this.VariantId
	return models.ContentRuleRowSave{
		Name:      this.Name(),
		VariantId: &id,
	}
}

// smallestVariantKey returns the lowest key in the map, so variant resolution
// is deterministic regardless of Go map iteration order.
func smallestVariantKey(variants map[int]string) int {
	if len(variants) == 0 {
		return 0
	}
	keys := make([]int, 0, len(variants))
	for key := range variants {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys[0]
}
