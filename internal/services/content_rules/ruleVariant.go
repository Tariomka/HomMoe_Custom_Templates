package content_rules

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
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
	VariantID int
}

// NewRuleVariant creates a variant rule. When no mapping is supplied it falls
// back to the dragon-utopia mapping (matching the C# dummy default). When no
// variant id is supplied it uses the lowest defined variant id for determinism.
// It returns an error when the resolved id is not present in the mapping.
func NewRuleVariant(mapping *models.VariantMapping, variantID *int) (*RuleVariant, error) {
	resolved := NewVariantMappingCatalog().GetDefaultMapping()
	if mapping != nil {
		resolved = *mapping
	}

	var id int
	if variantID != nil {
		id = *variantID
	} else {
		variantIDs := resolved.GetVariantIDsInOrder()
		if len(variantIDs) == 0 {
			return nil, common_errors.ErrNoVariantProvided
		}
		id = variantIDs[0]
	}

	if _, ok := resolved.GetVariantByID(id); !ok {
		return nil, fmt.Errorf("selected variant id %d is not present in the provided variant mapping", id)
	}

	return &RuleVariant{Mapping: resolved, VariantID: id}, nil
}

func (this *RuleVariant) Name() string { return RuleVariantName }

func (this *RuleVariant) Description() string { return RuleVariantDescription }

func (this *RuleVariant) Marker() string { return RuleVariantMarker }

func (this *RuleVariant) DisplayText() string {
	if description, ok := this.Mapping.GetVariantByID(this.VariantID); ok {
		return fmt.Sprintf("%s: %s", this.Name(), description)
	}

	return fmt.Sprintf("%s: Unforeseen Error", this.Name())
}

func (this *RuleVariant) Apply(item *entities.MandatoryContentItem) {
	id := this.VariantID
	item.Variant = &id
}

func (this *RuleVariant) SerializeToRowSave() models.ContentRuleRow {
	id := this.VariantID
	return models.ContentRuleRow{
		Name:      this.Name(),
		VariantID: &id,
	}
}
