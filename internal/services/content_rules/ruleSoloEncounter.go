package content_rules

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// Rule metadata constants for the solo-encounter rule.
const (
	RuleSoloEncounterName        = "Solo Encounter"
	RuleSoloEncounterDescription = "Solo encounter means that the content item will be spawned without any additional content items around, enforcing consistent guard strength. Setting to false will make it more likely to be spawned with other content items, but will not always guarantee it."
	RuleSoloEncounterMarker      = "S"
)

// RuleSoloEncounter forces the content item to spawn as a solo encounter.
type RuleSoloEncounter struct {
	IsSoloEncounter bool
}

// NewRuleSoloEncounter creates a solo-encounter rule with the supplied state.
func NewRuleSoloEncounter(isSoloEncounter bool) *RuleSoloEncounter {
	return &RuleSoloEncounter{IsSoloEncounter: isSoloEncounter}
}

func (this *RuleSoloEncounter) Name() string { return RuleSoloEncounterName }

func (this *RuleSoloEncounter) Description() string { return RuleSoloEncounterDescription }

// Marker shows "S" when solo and "!S" when explicitly not solo.
func (this *RuleSoloEncounter) Marker() string {
	if this.IsSoloEncounter {
		return RuleSoloEncounterMarker
	}

	return "!" + RuleSoloEncounterMarker
}

func (this *RuleSoloEncounter) DisplayText() string {
	return fmt.Sprintf("%s: %t", this.Name(), this.IsSoloEncounter)
}

func (this *RuleSoloEncounter) Apply(item *entities.MandatoryContentItem) {
	item.SoloEncounter = this.IsSoloEncounter
}

func (this *RuleSoloEncounter) SerializeToRowSave() models.ContentRuleRow {
	value := this.IsSoloEncounter
	return models.ContentRuleRow{
		Name:            this.Name(),
		IsSoloEncounter: &value,
	}
}
