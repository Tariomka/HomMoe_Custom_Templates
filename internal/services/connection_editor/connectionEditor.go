// Package connection_editor contains the model- and logic-layer behaviour of the
// Zone Connection Editor: guard-strength presets, zone tier classification,
// connection cloning, and graph-validation helpers. The visual canvas itself is
// rendered by the GUI layer; everything testable lives here.
package connection_editor

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// ZoneTier classifies a zone by its expected guard-strength bracket.
type ZoneTier int

const (
	ZoneTierBronze ZoneTier = iota
	ZoneTierSilver
	ZoneTierGold
	ZoneTierPlayerToPlayer
)

// StrengthLabels are the column labels of the guard-preset table.
var StrengthLabels = []string{"Weak", "Moderate", "Medium", "High", "Very High"}

// GuardPresets holds the guard values indexed by [tier][strength], matching the
// C# ZoneConnectionEditorWindow.GuardPresets table.
var GuardPresets = [4][5]int{
	{3000, 6000, 9000, 12000, 16000},    // Bronze
	{18000, 21000, 24000, 27000, 30000}, // Silver
	{36000, 42000, 48000, 54000, 60000}, // Gold
	{10000, 22000, 34000, 46000, 58000}, // PlayerToPlayer
}

// GuardPresetExtra is a named guard value shown at the top of the dropdown.
type GuardPresetExtra struct {
	Label string
	Value int
}

// TierExtras holds the "Generator Default" value for each tier.
var TierExtras = [4][]GuardPresetExtra{
	{{Label: "Generator Default", Value: 15000}}, // Bronze
	{{Label: "Generator Default", Value: 20000}}, // Silver
	{{Label: "Generator Default", Value: 25000}}, // Gold
	{{Label: "Generator Default", Value: 30000}}, // PlayerToPlayer
}

// WeeklyIncrementLabels and WeeklyIncrementValues describe the guard-growth presets.
var (
	WeeklyIncrementLabels = []string{"Slow (5%)", "Normal (10%)", "Standard (15%)", "Fast (20%)", "Very Fast (25%)"}
	WeeklyIncrementValues = []float64{0.05, 0.10, 0.15, 0.20, 0.25}
)

// UserCreatableConnectionTypes lists the connection types a user may add in the
// editor. Proximity is generator-only and intentionally excluded.
func UserCreatableConnectionTypes() []string {
	return []string{"Direct", "Portal"}
}

// GuardPresetsForTier returns the five guard-strength values for the given tier.
func GuardPresetsForTier(tier ZoneTier) [5]int {
	return GuardPresets[int(tier)]
}

// ExtrasForTier returns the named extra guard values for the given tier.
func ExtrasForTier(tier ZoneTier) []GuardPresetExtra {
	return TierExtras[int(tier)]
}

// GetZoneTier classifies a zone by name, using its guarded content pool to decide
// the bracket for neutral zones.
func GetZoneTier(zoneName string, zones []entities.Zone, playerZoneNames map[string]bool) ZoneTier {
	if zoneName == "" {
		return ZoneTierBronze
	}

	var zone *entities.Zone
	for i := range zones {
		if zones[i].Name == zoneName {
			zone = &zones[i]
			break
		}
	}
	if zone == nil {
		return ZoneTierBronze
	}

	if playerZoneNames[zone.Name] {
		return ZoneTierBronze
	}
	if zone.Name == "Hub" || strings.HasPrefix(zone.Name, "Hub-") {
		return ZoneTierBronze
	}

	if strings.HasPrefix(zone.Name, "Neutral-") {
		pool := ""
		if len(zone.GuardedContentPool) > 0 {
			pool = zone.GuardedContentPool[0]
		}
		if strings.Contains(pool, "_t4_") || strings.Contains(pool, "_t5_") {
			return ZoneTierGold
		}
		if strings.Contains(pool, "_t1_") || strings.Contains(pool, "_t2_") {
			return ZoneTierBronze
		}
		return ZoneTierSilver
	}

	return ZoneTierBronze
}

// HigherTierOf returns the higher guard tier of two zones, or PlayerToPlayer when
// both endpoints are player zones.
func HigherTierOf(zoneA, zoneB string, zones []entities.Zone, playerZoneNames map[string]bool) ZoneTier {
	aIsPlayer := zoneA != "" && playerZoneNames[zoneA]
	bIsPlayer := zoneB != "" && playerZoneNames[zoneB]
	if aIsPlayer && bIsPlayer {
		return ZoneTierPlayerToPlayer
	}

	tierA := GetZoneTier(zoneA, zones, playerZoneNames)
	tierB := GetZoneTier(zoneB, zones, playerZoneNames)
	if tierA > tierB {
		return tierA
	}
	return tierB
}

// CloneConnection returns a copy of c with IsUserAdded set as requested. Like the
// C# CloneConnection it shares the placement-rule slices with the original (a
// shallow reference copy), which is all the editor needs.
func CloneConnection(c entities.Connection, isUserAdded bool) entities.Connection {
	clone := c
	clone.IsUserAdded = isUserAdded
	return clone
}

// NewDefaultConnection builds a user-added Direct connection between two zones,
// seeded with the tier's generator-default guard value and the standard (15%)
// weekly increment, mirroring the C# AddConnectionWithDefaults.
func NewDefaultConnection(from, to string, zones []entities.Zone, playerZoneNames map[string]bool) entities.Connection {
	tier := HigherTierOf(from, to, zones, playerZoneNames)
	return entities.Connection{
		From:                 from,
		To:                   to,
		ConnectionType:       "Direct",
		GuardValue:           TierExtras[int(tier)][0].Value,
		GuardZone:            from,
		GuardMatchGroup:      "rnd_guard_" + ZoneLetterFromName(from) + "_" + ZoneLetterFromName(to),
		GuardWeeklyIncrement: WeeklyIncrementValues[2],
		IsUserAdded:          true,
	}
}

// ZoneLetterFromName extracts the identifier after the first dash in a zone name,
// e.g. "Spawn-A" → "A", "Neutral-C" → "C".
func ZoneLetterFromName(zoneName string) string {
	if _, after, ok := strings.Cut(zoneName, "-"); ok {
		return after
	}
	return zoneName
}

// FindIsolatedZones returns the names of zones not referenced by any connection.
func FindIsolatedZones(zones []entities.Zone, connections []entities.Connection) []string {
	var isolated []string
	for _, zone := range zones {
		referenced := false
		for _, connection := range connections {
			if connection.From == zone.Name || connection.To == zone.Name {
				referenced = true
				break
			}
		}
		if !referenced {
			isolated = append(isolated, zone.Name)
		}
	}
	return isolated
}

// ComputeHasErrors reports whether any connection references a zone name that
// does not exist in the zone list.
func ComputeHasErrors(zones []entities.Zone, connections []entities.Connection) bool {
	zoneNames := make(map[string]bool, len(zones))
	for _, zone := range zones {
		zoneNames[zone.Name] = true
	}
	for _, connection := range connections {
		if !zoneNames[connection.From] || !zoneNames[connection.To] {
			return true
		}
	}
	return false
}

// HasDuplicateName reports whether a connection other than current shares its
// (case-insensitive) name. current must be an element of connections; identity is
// compared by pointer to mirror the C# ReferenceEquals check.
func HasDuplicateName(connections []entities.Connection, current *entities.Connection) bool {
	if current == nil || len(current.Name) == 0 {
		return false
	}
	for i := range connections {
		if &connections[i] == current {
			continue
		}
		if strings.EqualFold(connections[i].Name, current.Name) {
			return true
		}
	}
	return false
}
