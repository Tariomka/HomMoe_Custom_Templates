package topology

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type RandomTopologyService struct {
	topologyBase
}

func NewRandomTopologyService() *RandomTopologyService {
	return &RandomTopologyService{
		topologyBase: newTopologyBase(),
	}
}

func (this *RandomTopologyService) GetTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	neutralLetters := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLetters[i] = nz.Label
	}
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1

	var allLabels []string
	if configuration.Topology == config.TopologyBalanced {
		allLabels = this.zoneLabelProvider.CreateBalancedRingZoneLabels(playerLabels, neutralZones, 0)
	} else {
		allLabels = append(append([]string{}, playerLabels...), neutralLetters...)
		rand.Shuffle(len(allLabels), func(i, j int) { allLabels[i], allLabels[j] = allLabels[j], allLabels[i] })
	}
	count := len(allLabels)

	var positions models.Positions
	if configuration.Topology == config.TopologyBalanced {
		positions = models.CreatePositionsFromPlans(allLabels, playerLabels, neutralZones)
	} else {
		for i := 0; i < count; i++ {
			positions.Add(models.NewPosition(rand.Float64()*0.9+0.05, rand.Float64()*0.9+0.05))
		}
	}

	pairs := positions.CreateDelaunayTriangulation()

	if configuration.Topology == config.TopologyBalanced {
		presentTiers := map[int]bool{}
		for _, label := range allLabels {
			tier := 0
			if !slices.Contains(playerLabels, label) {
				tier = neutralZones.GetTier(label)
			}
			presentTiers[tier] = true
		}
		var filtered [][2]int
		for _, pair := range pairs {
			tierA := 0
			if !slices.Contains(playerLabels, allLabels[pair[0]]) {
				tierA = neutralZones.GetTier(allLabels[pair[0]])
			}
			tierB := 0
			if !slices.Contains(playerLabels, allLabels[pair[1]]) {
				tierB = neutralZones.GetTier(allLabels[pair[1]])
			}
			low, high := tierA, tierB
			if low > high {
				low, high = high, low
			}
			if high-low <= 1 {
				filtered = append(filtered, pair)
				continue
			}
			skip := false
			for t := low + 1; t < high; t++ {
				if presentTiers[t] {
					skip = true
					break
				}
			}
			if !skip {
				filtered = append(filtered, pair)
			}
		}
		pairs = filtered
	}

	connsByZone := make(map[int][]string, count)
	var conns []template.Connection
	for _, p := range pairs {
		a, b := p[0], p[1]
		labelFrom := allLabels[a]
		labelTo := allLabels[b]
		if isIsolated && slices.Contains(playerLabels, labelFrom) && slices.Contains(playerLabels, labelTo) {
			continue
		}
		cn := fmt.Sprintf("Rnd-%s-%s", labelFrom, labelTo)
		connsByZone[a] = append(connsByZone[a], cn)
		connsByZone[b] = append(connsByZone[b], cn)
		conns = append(conns, template.Connection{
			Name:           cn,
			From:           this.zoneLabelProvider.CreateZoneName(labelFrom, playerLabels),
			To:             this.zoneLabelProvider.CreateZoneName(labelTo, playerLabels),
			ConnectionType: "Direct", GuardZone: this.zoneLabelProvider.CreateZoneName(labelFrom, playerLabels), SimTurnSquad: true,
			GuardValue: this.GetBorderGuardValue(labelFrom, labelTo, playerLabels, neutralZones, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("rnd_guard_%s_%s", labelFrom, labelTo),
		})
	}

	var zones []template.Zone
	for i, label := range allLabels {
		myConns := connsByZone[i]
		if pi := slices.Index(playerLabels, label); pi >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", pi+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles,
					configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, label == holdCityNeutralLetter))
		}
	}

	// Stamp generator-driven positions onto the freshly built zones so the
	// preview renderer can reproduce the exact geometry used to derive the
	// Delaunay connections. Balanced layouts also stamp the concentric ring
	// index so the preview can snap zones to clean rings
	for i := range zones {
		p := positions[i]
		zones[i].GeneratorPosition = &[2]float64{p.X, p.Y}
		if configuration.Topology == config.TopologyBalanced {
			tier := 0
			if !slices.Contains(playerLabels, allLabels[i]) {
				tier = neutralZones.GetTier(allLabels[i])
			}
			zones[i].GeneratorRing = &tier
		}
	}

	if configuration.RandomPortals {
		conns = append(conns, this.CreateRandomPortalConnections(playerLabels, allLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isIsolated {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLabels, zones, conns, tuning)...)
	}
	conns = this.CreateMissingConnections(playerLabels, allLabels, positions, zones, conns, tuning, neutralZones)
	return this.CreateVariant(playerLabels, allLabels[0], count, zones, conns)
}
