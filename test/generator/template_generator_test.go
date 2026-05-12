package generator

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

func TestGenerateUsesRequestedSettings(t *testing.T) {
	settings := &models.GeneratorSettings{
		TemplateName:    "Test Template",
		GameMode:        "Classic",
		PlayerCount:     4,
		MapSize:         "L",
		Topology:        models.TopologyDefault,
		ShowDescription: true,
		AdvancedSettings: &models.AdvancedSettings{
			GuardRandomization:     0.05,
			ConnectionCountPerZone: 2,
		},
	}

	template, err := generator.Generate(settings)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if template.Name != settings.TemplateName {
		t.Errorf("Expected name %s, got %s", settings.TemplateName, template.Name)
	}

	if template.Size != settings.MapSize {
		t.Errorf("Expected size %s, got %s", settings.MapSize, template.Size)
	}
}

func TestGenerateDefaultTopologyCreatesRingConnections(t *testing.T) {
	settings := &models.GeneratorSettings{
		TemplateName: "Ring Test",
		PlayerCount:  4,
		MapSize:      "L",
		Topology:     models.TopologyDefault,
		AdvancedSettings: &models.AdvancedSettings{
			GuardRandomization:     0.05,
			ConnectionCountPerZone: 2,
		},
	}

	template, err := generator.Generate(settings)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	variant := template.Variants[0]
	if len(variant.Connections) == 0 {
		t.Errorf("Expected connections, got none")
	}
}

func TestComputeContentScale(t *testing.T) {
	tests := []struct {
		mapSize   string
		zoneCount int
	}{
		{"S", 4},
		{"L", 4},
		{"2XL", 16},
	}

	for _, test := range tests {
		scale := generator.ComputeContentScale(test.mapSize, test.zoneCount)
		if scale < generator.ContentScaleMin || scale > generator.ContentScaleMax {
			t.Errorf("Scale out of bounds for %s: %f", test.mapSize, scale)
		}
	}
}

func TestGenerateWithRoadsDisabled(t *testing.T) {
	settings := &models.GeneratorSettings{
		TemplateName: "No Roads Test",
		PlayerCount:  4,
		MapSize:      "L",
		Topology:     models.TopologyDefault,
		AllowRoads:   false,
		AdvancedSettings: &models.AdvancedSettings{
			GuardRandomization:     0.05,
			ConnectionCountPerZone: 2,
		},
	}

	template, err := generator.Generate(settings)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	variant := template.Variants[0]
	for _, conn := range variant.Connections {
		if conn.HasRoad {
			t.Errorf("Expected no roads, but connection %s-%s has roads", conn.FromZone, conn.ToZone)
		}
	}
}
