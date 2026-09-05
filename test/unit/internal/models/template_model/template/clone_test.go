package template_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateIsCloned_CloneEqualsTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	source := newAllFieldsModel(t)

	// Act
	clone := source.Clone()

	// Assert
	assert.Equal(t, source, clone)
}

func TestWhenTheCloneIsMutated_SourceIsUntouched(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		mutate func(clone *template_model.Template)
	}{
		{"WhenAValueOverrideChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.ValueOverrides[0].GuardValue++
		}},
		{"WhenTheOrientationChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Orientation.Mode = "mutated"
		}},
		{"WhenABorderNoiseChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Border.ObstaclesNoise[0].Frequency++
		}},
		{"WhenABonusParameterChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.GameRules.Bonuses[0].Parameters[0] = "mutated"
		}},
		{"WhenATournamentDayChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.GameRules.WinConditions.TournamentDays[0]++
		}},
		{"WhenAGameRuleBanChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.GameRules.GlobalBans.Items[0] = "mutated"
		}},
		{"WhenATemplateBanChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.GlobalBans.Magics[0] = "mutated"
		}},
		{"WhenAZoneTierChanges_SourceKeepsIts", func(clone *template_model.Template) {
			*clone.Variants[0].Zones[0].Quality = neutral_zone.QualityLowest
		}},
		{"WhenAZonePoolChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Variants[0].Zones[0].GuardedContentPool[0] = "mutated"
		}},
		{"WhenAZoneBiomeArgChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Variants[0].Zones[0].ZoneBiome.Args[0] = "mutated"
		}},
		{"WhenAMainObjectFactionChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Variants[0].Zones[0].MainObjects[0].Faction.Type = "mutated"
		}},
		{"WhenARoadEndpointChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Variants[0].Zones[0].Roads[0].From.Args[0] = "mutated"
		}},
		{"WhenAConnectionFlagChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Variants[0].Connections[0].IsUserAdded = false
		}},
		{"WhenAPortalRuleArgChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.Variants[0].Connections[0].PortalPlacementRulesFrom[0].Args[0] = "mutated"
		}},
		{"WhenAZoneLayoutFractionChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.ZoneLayouts[0].GuardedEncounterResourceFractions.Fractions[0]++
		}},
		{"WhenAMandatoryItemRuleArgChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.MandatoryContent[0].Content[0].Rules[0].Args[0] = "mutated"
		}},
		{"WhenAContentLimitListChanges_SourceKeepsIts", func(clone *template_model.Template) {
			clone.ContentCountLimits[0].Limits[0].IncludeLists[0] = "mutated"
		}},
		{"WhenAContentPoolEntryChanges_SourceKeepsIts", func(clone *template_model.Template) {
			for key := range clone.ContentPools[0] {
				clone.ContentPools[0][key] = "mutated"
			}
		}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			source := newAllFieldsModel(t)
			pristine := newAllFieldsModel(t)
			clone := source.Clone()

			// Act
			testCase.mutate(&clone)

			// Assert
			assert.Equal(t, pristine, source)
		})
	}
}

func TestWhenSlicesAreNil_CloneKeepsThemNil(t *testing.T) {
	t.Parallel()
	// Arrange
	source := template_model.Template{Variants: []template_model.Variant{{}}}

	// Act
	clone := source.Clone()

	// Assert
	assert.Equal(t, source, clone)
}

func TestWhenSlicesAreEmpty_CloneKeepsThemEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	source := template_model.Template{
		ContentPools: []template_model.ContentPool{},
		ContentLists: []template_model.ContentList{},
	}

	// Act
	clone := source.Clone()

	// Assert
	assert.Equal(t, source, clone)
}

// newAllFieldsModel lifts the fuzzed all-fields fixture and stamps the two
// model-only fields the entity cannot carry, so the clone is proved on the whole
// model and not just on the wire-format subset.
func newAllFieldsModel(t *testing.T) template_model.Template {
	t.Helper()
	model := template_model.ToTemplateModel(test_helpers.NewAllFieldsTemplate())
	require.NotEmpty(t, model.Variants)
	require.NotEmpty(t, model.Variants[0].Zones)
	require.NotEmpty(t, model.Variants[0].Connections)
	model.Variants[0].Zones[0].Quality = new(neutral_zone.QualityHigh)
	model.Variants[0].Connections[0].IsUserAdded = true
	return model
}
