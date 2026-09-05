package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
)

// NewAllFieldsTemplate builds a .rmg.json template with every field populated,
// so a round trip through the template mapper proves the mapper carries the
// whole schema. Fuzzing it rather than hand-writing it is the point: a field
// added to the schema and forgotten in a converter fails the round trip instead
// of passing silently, which a hand-written fixture could never catch.
func NewAllFieldsTemplate() entities.RmgTemplate {
	faker := gofakeit.New(allFieldsTemplateSeed)

	var template entities.RmgTemplate
	if err := faker.Struct(&template); err != nil {
		panic(err)
	}

	// gofakeit cannot invent a value for an `any`, so the loosely-typed corners
	// of the schema are filled by hand.
	fillLooselyTypedFields(&template, faker)

	return template
}

// allFieldsTemplateSeed keeps the fixture deterministic across runs, so a
// failure is always reproducible.
const allFieldsTemplateSeed = 20260902

func fillLooselyTypedFields(template *entities.RmgTemplate, faker *gofakeit.Faker) {
	for index := range template.ContentPools {
		template.ContentPools[index] = entities.ContentPool{faker.Word(): faker.Word()}
	}
	for index := range template.ContentLists {
		template.ContentLists[index] = entities.ContentList{faker.Word(): faker.Word()}
	}

	for contentIndex := range template.MandatoryContent {
		items := template.MandatoryContent[contentIndex].Content
		for itemIndex := range items {
			fillPlacementRuleArgs(items[itemIndex].Rules, faker)
		}
	}

	for variantIndex := range template.Variants {
		connections := template.Variants[variantIndex].Connections
		for connectionIndex := range connections {
			fillPlacementRuleArgs(connections[connectionIndex].PortalPlacementRulesFrom, faker)
			fillPlacementRuleArgs(connections[connectionIndex].PortalPlacementRulesTo, faker)
		}
	}
}

func fillPlacementRuleArgs(rules []entities.PlacementRule, faker *gofakeit.Faker) {
	for index := range rules {
		rules[index].Args = []any{faker.Word(), faker.Number(1, 100)}
	}
}
