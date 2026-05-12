package services

import (
	"fmt"
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

const (
	DefaultGuardRandomization = 0.05
	DefaultConnectionPerZone  = 2
	MaxZones                  = 32
	ContentScaleMin           = 0.5
	ContentScaleMax           = 2.5
)

var zoneLetters = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P",
	"Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF",
}

type NeutralZonePlan struct {
	Letter      string
	Quality     models.NeutralZoneQuality
	CastleCount int
}

type TopologyAdjacency struct {
	Connections map[string][]string
}

// Generate creates an RmgTemplate from GeneratorSettings
func Generate(settings *models.GeneratorSettings) (*models.RmgTemplate, error) {
	if settings == nil {
		return nil, fmt.Errorf("settings cannot be nil")
	}

	if settings.AdvancedSettings == nil {
		settings.AdvancedSettings = &models.AdvancedSettings{
			GuardRandomization:     DefaultGuardRandomization,
			ConnectionCountPerZone: DefaultConnectionPerZone,
		}
	}

	neutralZonePlan := buildNeutralZonePlan(settings)
	holdCityLetter := ""
	if settings.EnableCityHold && len(neutralZonePlan) > 0 {
		holdCityLetter = pickHoldCityNeutralLetter(settings.PlayerCount, neutralZonePlan)
	}

	adjacency := buildTopologyAdjacency(settings, neutralZonePlan)
	gameRules := buildGameRules(settings, holdCityLetter)
	variant := buildVariant(settings, neutralZonePlan, adjacency)

	template := &models.RmgTemplate{
		Name:        settings.TemplateName,
		Description: buildDescription(settings),
		Size:        settings.MapSize,
		GameRules:   gameRules,
		Variants:    []models.Variant{*variant},
	}

	return template, nil
}

func buildNeutralZonePlan(settings *models.GeneratorSettings) []*NeutralZonePlan {
	var plan []*NeutralZonePlan

	advancedSettings := settings.AdvancedSettings
	if advancedSettings == nil {
		advancedSettings = &models.AdvancedSettings{
			NeutralZoneLowCount:    1,
			NeutralZoneMediumCount: 1,
			NeutralZoneHighCount:   0,
		}
	}

	for i := 0; i < advancedSettings.NeutralZoneLowCount; i++ {
		plan = append(plan, &NeutralZonePlan{
			Letter:      getZoneLetter(settings.PlayerCount + len(plan)),
			Quality:     models.QualityLow,
			CastleCount: 0,
		})
	}

	for i := 0; i < advancedSettings.NeutralZoneMediumCount; i++ {
		plan = append(plan, &NeutralZonePlan{
			Letter:      getZoneLetter(settings.PlayerCount + len(plan)),
			Quality:     models.QualityMedium,
			CastleCount: 1,
		})
	}

	for i := 0; i < advancedSettings.NeutralZoneHighCount; i++ {
		plan = append(plan, &NeutralZonePlan{
			Letter:      getZoneLetter(settings.PlayerCount + len(plan)),
			Quality:     models.QualityHigh,
			CastleCount: 1,
		})
	}

	return plan
}

func getZoneLetter(i int) string {
	if i < 0 || i >= len(zoneLetters) {
		return ""
	}
	return zoneLetters[i]
}

func pickHoldCityNeutralLetter(playerCount int, neutralZonePlan []*NeutralZonePlan) string {
	if len(neutralZonePlan) == 0 {
		return ""
	}
	return neutralZonePlan[0].Letter
}

func buildTopologyAdjacency(settings *models.GeneratorSettings, neutralZonePlan []*NeutralZonePlan) *TopologyAdjacency {
	adjacency := &TopologyAdjacency{
		Connections: make(map[string][]string),
	}

	playerLetters := make([]string, settings.PlayerCount)
	for i := 0; i < settings.PlayerCount; i++ {
		playerLetters[i] = getZoneLetter(i)
	}

	neutralLetters := make([]string, len(neutralZonePlan))
	for i, plan := range neutralZonePlan {
		neutralLetters[i] = plan.Letter
	}

	switch settings.Topology {
	case models.TopologyDefault:
		// Ring topology
		for i, letter := range playerLetters {
			next := playerLetters[(i+1)%len(playerLetters)]
			adjacency.Connections[letter] = append(adjacency.Connections[letter], next)
		}
		for i, letter := range playerLetters {
			for j, neutral := range neutralLetters {
				if (i+j)%2 == 0 {
					adjacency.Connections[letter] = append(adjacency.Connections[letter], neutral)
					adjacency.Connections[neutral] = append(adjacency.Connections[neutral], letter)
				}
			}
		}

	case models.TopologyHubAndSpoke:
		hub := getZoneLetter(settings.PlayerCount)
		if len(neutralLetters) > 0 {
			hub = neutralLetters[0]
		}
		for _, letter := range playerLetters {
			adjacency.Connections[letter] = append(adjacency.Connections[letter], hub)
			adjacency.Connections[hub] = append(adjacency.Connections[hub], letter)
		}

	case models.TopologyChain:
		allZones := append(append([]string{}, playerLetters...), neutralLetters...)
		for i := 0; i < len(allZones)-1; i++ {
			adjacency.Connections[allZones[i]] = append(adjacency.Connections[allZones[i]], allZones[i+1])
			adjacency.Connections[allZones[i+1]] = append(adjacency.Connections[allZones[i+1]], allZones[i])
		}

	case models.TopologySharedWeb:
		for i, player := range playerLetters {
			neutral := neutralLetters[i%len(neutralLetters)]
			nextPlayer := playerLetters[(i+1)%len(playerLetters)]
			adjacency.Connections[player] = append(adjacency.Connections[player], neutral)
			adjacency.Connections[neutral] = append(adjacency.Connections[neutral], nextPlayer)
		}

	case models.TopologyRandom:
		for i, letter := range playerLetters {
			next := playerLetters[(i+1)%len(playerLetters)]
			adjacency.Connections[letter] = append(adjacency.Connections[letter], next)
		}
	}

	return adjacency
}

func buildGameRules(settings *models.GeneratorSettings, holdCityLetter string) models.GameRules {
	rules := models.GameRules{
		HeroCount: 1,
	}

	if settings.GameEndConditions != nil {
		if settings.GameEndConditions.EnableClassicVictory {
			rules.WinConditions = append(rules.WinConditions, models.WinCondition{
				Type:      "ClassicVictory",
				Condition: "captureAllTowns",
			})
		}
		if settings.GameEndConditions.EnableCityHold && holdCityLetter != "" {
			rules.WinConditions = append(rules.WinConditions, models.WinCondition{
				Type: "CityHold",
				Parameters: map[string]interface{}{
					"town": holdCityLetter,
					"days": 7,
				},
			})
		}
	} else {
		rules.WinConditions = append(rules.WinConditions, models.WinCondition{
			Type:      "ClassicVictory",
			Condition: "captureAllTowns",
		})
	}

	return rules
}

func buildVariant(settings *models.GeneratorSettings, neutralZonePlan []*NeutralZonePlan, adjacency *TopologyAdjacency) *models.Variant {
	variant := &models.Variant{
		Name:  "Default",
		Zones: buildZones(settings, neutralZonePlan),
	}

	for zone := range adjacency.Connections {
		for _, targetZone := range adjacency.Connections[zone] {
			if zone < targetZone {
				variant.Connections = append(variant.Connections, models.Connection{
					FromZone: zone,
					ToZone:   targetZone,
					HasRoad:  settings.AllowRoads,
					PortalPlacement: models.PortalPlacement{
						From: models.PortalEndpoint{ZoneName: zone},
						To:   models.PortalEndpoint{ZoneName: targetZone},
					},
				})
			}
		}
	}

	return variant
}

func buildZones(settings *models.GeneratorSettings, neutralZonePlan []*NeutralZonePlan) []models.Zone {
	var zones []models.Zone

	for i := 0; i < settings.PlayerCount; i++ {
		zone := models.Zone{
			Name:   fmt.Sprintf("Player %d", i+1),
			Type:   "player",
			Letter: getZoneLetter(i),
			Owner:  i,
			Layout: models.ZoneLayout{Type: "default"},
			GuardSettings: models.GuardSettings{
				Randomization: settings.AdvancedSettings.GuardRandomization,
			},
			ContentPools: models.ContentPools{},
		}
		zones = append(zones, zone)
	}

	for _, plan := range neutralZonePlan {
		zone := models.Zone{
			Name:   fmt.Sprintf("Neutral %s (%s)", plan.Letter, plan.Quality),
			Type:   "neutral",
			Letter: plan.Letter,
			Layout: models.ZoneLayout{Type: "default"},
			GuardSettings: models.GuardSettings{
				Randomization: settings.AdvancedSettings.GuardRandomization,
			},
			ContentPools: models.ContentPools{},
		}
		zones = append(zones, zone)
	}

	return zones
}

func buildDescription(settings *models.GeneratorSettings) string {
	var parts []string
	parts = append(parts, settings.TemplateName)

	if settings.ShowDescription {
		parts = append(parts, fmt.Sprintf("Players: %d", settings.PlayerCount))
		parts = append(parts, fmt.Sprintf("Size: %s", settings.MapSize))
		parts = append(parts, fmt.Sprintf("Mode: %s", settings.GameMode))
		parts = append(parts, fmt.Sprintf("Topology: %s", settings.Topology))

		if settings.AllowFootholds {
			parts = append(parts, "Includes Footholds")
		}
	}

	return strings.Join(parts, " | ")
}

// ComputeContentScale calculates content scaling
func ComputeContentScale(mapSize string, zoneCount int) float64 {
	sizeScale := 1.0
	switch mapSize {
	case "S":
		sizeScale = 0.7
	case "M":
		sizeScale = 0.85
	case "L":
		sizeScale = 1.0
	case "XL":
		sizeScale = 1.2
	case "2XL":
		sizeScale = 1.4
	}

	zoneAdjustment := math.Sqrt(float64(zoneCount) / 4)
	scale := sizeScale * zoneAdjustment

	if scale < ContentScaleMin {
		return ContentScaleMin
	}
	if scale > ContentScaleMax {
		return ContentScaleMax
	}
	return scale
}
