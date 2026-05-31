package topology

import (
	"fmt"
	"math/rand/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

func buildVariantRandom(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	neutralLetters := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLetters[i] = nz.Letter
	}
	isolate := settings.NoDirectPlayerConnections && len(playerLetters) > 1

	var allLetters []string
	if settings.Topology == config.TopologyBalanced {
		allLetters = buildBalancedRingLetters(playerLetters, neutralZones, 0)
	} else {
		allLetters = append(append([]string{}, playerLetters...), neutralLetters...)
		rand.Shuffle(len(allLetters), func(i, j int) { allLetters[i], allLetters[j] = allLetters[j], allLetters[i] })
	}
	count := len(allLetters)

	var pos [][2]float64
	if settings.Topology == config.TopologyBalanced {
		pos = buildBalancedRandomPositions(allLetters, playerLetters, neutralByLetter)
	} else {
		for i := 0; i < count; i++ {
			pos = append(pos, [2]float64{rand.Float64()*0.9 + 0.05, rand.Float64()*0.9 + 0.05})
		}
	}

	pairs := delaunayEdges(pos)

	if settings.Topology == config.TopologyBalanced {
		presentTiers := map[int]bool{}
		for _, l := range allLetters {
			presentTiers[zoneTierRank(l, playerLetters, neutralByLetter)] = true
		}
		var filtered [][2]int
		for _, p := range pairs {
			ta := zoneTierRank(allLetters[p[0]], playerLetters, neutralByLetter)
			tb := zoneTierRank(allLetters[p[1]], playerLetters, neutralByLetter)
			lo, hi := ta, tb
			if lo > hi {
				lo, hi = hi, lo
			}
			if hi-lo <= 1 {
				filtered = append(filtered, p)
				continue
			}
			skip := false
			for t := lo + 1; t < hi; t++ {
				if presentTiers[t] {
					skip = true
					break
				}
			}
			if !skip {
				filtered = append(filtered, p)
			}
		}
		pairs = filtered
	}

	connsByZone := make(map[int][]string, count)
	var conns []template.Connection
	for _, p := range pairs {
		a, b := p[0], p[1]
		from := allLetters[a]
		to := allLetters[b]
		if isolate && contains(playerLetters, from) && contains(playerLetters, to) {
			continue
		}
		cn := fmt.Sprintf("Rnd-%s-%s", from, to)
		connsByZone[a] = append(connsByZone[a], cn)
		connsByZone[b] = append(connsByZone[b], cn)
		conns = append(conns, template.Connection{
			Name: cn, From: zoneName(from, playerLetters), To: zoneName(to, playerLetters),
			ConnectionType: "Direct", GuardZone: zoneName(from, playerLetters), SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("rnd_guard_%s_%s", from, to),
		})
	}

	var zones []template.Zone
	for i := 0; i < count; i++ {
		letter := allLetters[i]
		myConns := connsByZone[i]
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, letter == holdCityNeutralLetter))
		}
	}

	// Stamp generator-driven positions onto the freshly built zones so the
	// preview renderer can reproduce the exact geometry used to derive the
	// Delaunay connections. Balanced layouts also stamp the concentric ring
	// index so the preview can snap zones to clean rings
	for i := range zones {
		p := pos[i]
		zones[i].GeneratorPosition = &[2]float64{p[0], p[1]}
		if settings.Topology == config.TopologyBalanced {
			r := zoneTierRank(allLetters[i], playerLetters, neutralByLetter)
			zones[i].GeneratorRing = &r
		}
	}

	if settings.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, allLetters, tuning, settings.MaxPortalConnections)...)
	}
	if isolate {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	ensureFullConnectivity(playerLetters, allLetters, pos, zones, &conns, tuning, neutralByLetter)
	return makeVariant(playerLetters, allLetters[0], count, zones, conns)
}
