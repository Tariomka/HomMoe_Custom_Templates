package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func main() {
	var (
		name      = flag.String("name", "Generated Template", "Template name")
		players   = flag.Int("players", 4, "Number of players")
		mapSize   = flag.String("size", "L", "Map size (S, M, L, XL, 2XL)")
		topology  = flag.String("topology", "Default", "Topology (Default, Chain, HubAndSpoke, SharedWeb, Random)")
		output    = flag.String("output", ".", "Output directory")
		gameMode  = flag.String("game", "Classic", "Game type (Classic, Blitz, Heroic)")
		roads     = flag.Bool("roads", true, "Allow roads")
		portals   = flag.Bool("portals", true, "Allow portals")
		footholds = flag.Bool("footholds", false, "Allow footholds")
		cityHold  = flag.Bool("cityhold", false, "Enable city hold win condition")
	)

	flag.Parse()

	// Create settings
	settings := &models.GeneratorSettings{
		TemplateName:                *name,
		GameMode:                    *gameMode,
		PlayerCount:                 *players,
		MapSize:                     *mapSize,
		Topology:                    models.MapTopology(*topology),
		AllowRoads:                  *roads,
		AllowPortals:                *portals,
		AllowFootholds:              *footholds,
		EnableCityHold:              *cityHold,
		ShowDescription:             true,
		IncludeOptionsInDescription: true,
		AdvancedSettings: &models.AdvancedSettings{
			GuardRandomization:     0.05,
			ConnectionCountPerZone: 2,
			NeutralZoneLowCount:    1,
			NeutralZoneMediumCount: 1,
			NeutralZoneHighCount:   0,
		},
		GameEndConditions: &models.GameEndConditions{
			EnableClassicVictory: true,
			EnableCityHold:       *cityHold,
		},
	}

	// Generate template
	template, err := services.Generate(settings)
	if err != nil {
		log.Fatalf("Failed to generate template: %v", err)
	}

	// Create filename
	filename := fmt.Sprintf("%s.rmg.json", settings.TemplateName)
	filepath := filepath.Join(*output, filename)

	// Marshal to JSON
	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal template: %v", err)
	}

	// Write to file
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("Template generated successfully: %s\n", filepath)
	fmt.Printf("  Name: %s\n", template.Name)
	fmt.Printf("  Description: %s\n", template.Description)
	fmt.Printf("  Size: %s\n", template.Size)
	fmt.Printf("  Zones: %d\n", len(template.Variants[0].Zones))
	fmt.Printf("  Connections: %d\n", len(template.Variants[0].Connections))
}
