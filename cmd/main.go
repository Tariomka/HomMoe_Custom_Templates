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
		mapSize   = flag.Int("size", 160, "Map size in tiles (e.g. 96, 128, 160, 192, 224)")
		topology  = flag.String("topology", "Default", "Topology (Default, Chain, HubAndSpoke, SharedWeb, Random)")
		output    = flag.String("output", ".", "Output directory")
		gameMode  = flag.String("game", "Classic", "Game type (Classic, Blitz, Heroic)")
		roads     = flag.Bool("roads", true, "Generate roads")
		portals   = flag.Bool("portals", false, "Random portals")
		footholds = flag.Bool("footholds", true, "Spawn remote footholds")
		cityHold  = flag.Bool("cityhold", false, "Enable city hold win condition")
		neutrals  = flag.Int("neutrals", 2, "Number of neutral zones")
	)

	flag.Parse()

	settings := models.NewGeneratorSettings()
	settings.TemplateName = *name
	settings.GameMode = *gameMode
	settings.PlayerCount = *players
	settings.MapSize = *mapSize
	settings.Topology = models.MapTopology(*topology)
	settings.GenerateRoads = *roads
	settings.RandomPortals = *portals
	settings.SpawnRemoteFootholds = *footholds
	settings.ZoneCfg.NeutralZoneCount = *neutrals
	if *cityHold {
		settings.GameEndConditions = &models.GameEndConditions{
			VictoryCondition: "win_condition_5",
			CityHold:         true,
			CityHoldDays:     6,
			LostStartCityDay: 3,
		}
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
	fmt.Printf("  Size: %dx%d\n", template.SizeX, template.SizeZ)
	fmt.Printf("  Zones: %d\n", len(template.Variants[0].Zones))
	fmt.Printf("  Connections: %d\n", len(template.Variants[0].Connections))
}
