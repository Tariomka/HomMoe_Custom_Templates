// Zone-mutation logic of the Manual Zone Editor: adding, deleting and
// re-profiling zones. The visual canvas lives in the GUI layer; everything
// testable lives here.
package connection_editor

import (
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

// zoneLabels mirrors the generator's label pool (zones.ZoneLabelProvider).
var zoneLabels = []string{
	"A", "B", "C", "D", "E", "F", "G", "H",
	"I", "J", "K", "L", "M", "N", "O", "P",
	"Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF",
}

// QualityLabels are the display names of the neutral-zone quality presets,
// indexed by models.NeutralZoneQuality.
var QualityLabels = []string{"Low", "Medium", "High"}

// NextFreeZoneLabel returns the first generator label not used by any zone, or
// "" when the pool is exhausted.
func NextFreeZoneLabel(zones []template.Zone) string {
	used := make(map[string]bool, len(zones))
	for _, zone := range zones {
		used[ZoneLetterFromName(zone.Name)] = true
	}
	for _, label := range zoneLabels {
		if !used[label] {
			return label
		}
	}
	return ""
}

// NewDefaultNeutralZone builds a manually-added neutral zone with the same
// builder the generator uses. The mandatory-content reference is cleared
// because no template-level definition exists for a manual zone.
func NewDefaultNeutralZone(
	label string,
	quality models.NeutralZoneQuality,
	castleCount int,
	generateRoads bool,
	tuning models.GenerationTuning) template.Zone {
	topology := base.NewTopologyBase()
	plan := models.NeutralZonePlan{Label: label, Quality: quality, CastleCount: castleCount}
	zone := topology.CreateNeutralZone(plan, nil, 1.0, false, generateRoads, tuning, false)
	zone.MandatoryContent = nil
	return zone
}

// QualityOfZone infers a neutral zone's quality preset from its guarded
// content pool, mirroring GetZoneTier's bracket rules.
func QualityOfZone(zone template.Zone) models.NeutralZoneQuality {
	pool := ""
	if len(zone.GuardedContentPool) > 0 {
		pool = zone.GuardedContentPool[0]
	}
	if strings.Contains(pool, "_t4_") || strings.Contains(pool, "_t5_") {
		return models.QualityHigh
	}
	if strings.Contains(pool, "_t1_") || strings.Contains(pool, "_t2_") {
		return models.QualityLow
	}
	return models.QualityMedium
}

// CountZoneCastles returns the number of City main objects in the zone.
func CountZoneCastles(zone template.Zone) int {
	count := 0
	for _, mainObject := range zone.MainObjects {
		if strings.EqualFold(mainObject.Type, "City") {
			count++
		}
	}
	return count
}

// ApplyNeutralZoneQuality re-applies the quality profile (layout, guard
// multiplier, content pools and values) and rebuilds the zone's castles for
// the requested count. Only meaningful for neutral zones.
func ApplyNeutralZoneQuality(
	zone *template.Zone,
	quality models.NeutralZoneQuality,
	castleCount int,
	tuning models.GenerationTuning) {
	profile := models.NewNeutralZoneProfile(quality)
	zone.Layout = profile.Layout
	zone.GuardMultiplier = tuning.ScaleByNeutralGuardStrengthPrecise(profile.GuardMultiplier)
	zone.GuardReactionDistribution = profile.GuardReactionDistribution
	zone.GuardedContentPool = profile.GuardedContentPool
	zone.UnguardedContentPool = profile.UnguardedContentPool
	zone.ResourcesContentPool = profile.ResourcesContentPool
	zone.GuardedContentValue = tuning.ScaleByStructureDensity(float64(profile.GuardedContentValue) * tuning.ContentScale)
	zone.GuardedContentValuePerArea = tuning.ScaleByStructureDensity(float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))
	zone.UnguardedContentValue = tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValue) * tuning.ContentScale)
	zone.UnguardedContentValuePerArea = tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))
	zone.ResourcesValue = tuning.ScaleByResourceDensity(float64(profile.ResourcesValue) * tuning.ContentScale)
	zone.ResourcesValuePerArea = tuning.ScaleByResourceDensity(float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale))
	zone.MainObjects = base.CreateNeutralZoneCastles(profile, tuning, castleCount, false)
}

// CanDeleteZone reports whether the zone may be removed in the editor. Spawn
// zones are owned by the General tab's player count and cannot be deleted.
func CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return !playerZoneNames[zoneName]
}

// RemoveZone returns the zone and connection lists without the named zone and
// without any connection referencing it.
func RemoveZone(
	zones []template.Zone,
	connections []template.Connection,
	zoneName string) ([]template.Zone, []template.Connection) {
	keptZones := make([]template.Zone, 0, len(zones))
	for _, zone := range zones {
		if zone.Name != zoneName {
			keptZones = append(keptZones, zone)
		}
	}
	keptConnections := make([]template.Connection, 0, len(connections))
	for _, connection := range connections {
		if connection.From != zoneName && connection.To != zoneName {
			keptConnections = append(keptConnections, connection)
		}
	}
	return keptZones, keptConnections
}

// FindOpenPosition returns a normalized position on a coarse interior grid
// that maximizes the distance to the occupied positions.
func FindOpenPosition(occupied [][2]float64) [2]float64 {
	const gridSteps = 7
	best := [2]float64{0.5, 0.5}
	bestScore := -1.0
	for row := range gridSteps {
		for col := range gridSteps {
			candidate := [2]float64{
				0.1 + 0.8*float64(col)/float64(gridSteps-1),
				0.1 + 0.8*float64(row)/float64(gridSteps-1),
			}
			score := math.MaxFloat64
			for _, p := range occupied {
				score = math.Min(score, math.Hypot(candidate[0]-p[0], candidate[1]-p[1]))
			}
			if len(occupied) == 0 {
				score = math.Hypot(candidate[0]-0.5, candidate[1]-0.5)
			}
			if score > bestScore {
				bestScore = score
				best = candidate
			}
		}
	}
	return best
}
