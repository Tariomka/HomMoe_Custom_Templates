package tournament_variant

import (
	"fmt"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

type RingClusterService struct {
	base.TopologyBase
}

func NewRingClusterService() *RingClusterService {
	return &RingClusterService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *RingClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans models.NeutralZonePlans,
	playerIndex int,
	playerLabel string) ([]template.Zone, []template.Connection) {
	var zones []template.Zone
	var connections []template.Connection
	sortedNeutrals := make([]neutralZonePlan, len(playerNeutralZonePlans))
	copy(sortedNeutrals, playerNeutralZonePlans)
	sort.SliceStable(sortedNeutrals, func(i, j int) bool {
		si, sj := neutralZoneBalanceScore(sortedNeutrals[i]), neutralZoneBalanceScore(sortedNeutrals[j])
		if si != sj {
			return si < sj
		}
		return sortedNeutrals[i].Letter < sortedNeutrals[j].Letter
	})

	n := len(sortedNeutrals)
	orderedNeutrals := make([]neutralZonePlan, n)
	lo, hi := 0, n-1
	for i := range n {
		if i%2 == 0 {
			orderedNeutrals[lo] = sortedNeutrals[i]
			lo++
		} else {
			orderedNeutrals[hi] = sortedNeutrals[i]
			hi--
		}
	}

	ringLetters := []string{playerLabel}
	for _, nz := range orderedNeutrals {
		ringLetters = append(ringLetters, nz.Letter)
	}
	count := len(ringLetters)
	if count < 2 {
		// Lone player zone — no ring edges possible.
		*zones = append(*zones, buildSpawnZone(playerLabel, fmt.Sprintf("Player%d", playerIndex+1), nil, configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		return
	}

	connNames := make([]string, count)
	for i := range count {
		next := (i + 1) % count
		connNames[i] = fmt.Sprintf("TRing-%s-%s", ringLetters[i], ringLetters[next])
	}

	for i, letter := range ringLetters {
		prev := (i - 1 + count) % count
		seen := map[string]bool{}
		var myConns []string
		for _, name := range []string{connNames[prev], connNames[i]} {
			if !seen[name] {
				seen[name] = true
				myConns = append(myConns, name)
			}
		}
		if i == 0 {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(allNeutralZonePlans[letter], myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, false))
		}
	}

	for i := range count {
		next := (i + 1) % count
		from := ringLetters[i]
		to := ringLetters[next]
		fromZone := "Spawn-" + from
		if i != 0 {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if next != 0 {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name: connNames[i], From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: GetBorderGuardValue(from, to, []string{playerLabel}, allNeutralZonePlans, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_ring_guard_%s_%s", from, to),
		})
	}

	return zones, connections
}
