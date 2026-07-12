package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
)

// FractalTopologyService grows one self-similar fractal per player. Every player
// zone is the base of its fractal and sits on the outer ring; its neutral zones
// branch inward toward the centre in nested tiers - low zones nearest the player
// fan out widely, then merge into fewer medium zones, then into the high zones
// that gather near the centre (the farthest point from any player). The innermost
// tips of neighbouring fractals are woven into a shared central ring, so the
// player sectors integrate seamlessly into one rotationally symmetric pattern and
// no two players ever border each other - the design itself keeps them apart,
// without relying on the NoDirectPlayerConnections flag.
type FractalTopologyService struct {
	RandomTopologyService
}

func NewFractalTopologyService() *FractalTopologyService {
	return &FractalTopologyService{
		RandomTopologyService: *NewRandomTopologyService(),
	}
}

func (this *FractalTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutralZone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	return this.createVariantFromLayout(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel, this.createFractalLayout)
}

// fractalTree holds the zone indices of one player's fractal. levels[0] are the
// low (tier 1) zones closest to the player, levels[1] the medium (tier 2) zones,
// levels[2] the high (tier 3) zones nearest the centre. Any level may be empty
// when the zone pool does not provide that tier for this player.
type fractalTree struct {
	levels [3][]int
	player int
}

// createFractalLayout builds the parallel label, position and connection-pair
// slices. Players are evenly spaced on the outer ring and each one anchors a
// fractal whose neutral tiers nest inward: the band of each tier narrows as the
// radius shrinks, so the wide spray of low zones funnels into the tight cluster
// of high zones at the centre - a self-similar, converging branch per player.
func (this *FractalTopologyService) createFractalLayout(
	playerLabels []string,
	neutralZones neutralZone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
	const (
		centreX      = 0.5
		centreY      = 0.5
		playerRadius = 0.45 // base of every fractal, on the outer ring
		startAngle   = -math.Pi / 2.0
	)
	// Radius of each neutral tier measured from the centre. Low sits just inside
	// the player, high gathers near the middle (the farthest point from a player).
	tierRadius := [3]float64{0.32, 0.19, 0.08}
	// Angular half-width of each tier's band as a fraction of the player's sector.
	// It shrinks inward so the branches visibly converge toward the centre.
	tierSpread := [3]float64{0.85, 0.52, 0.26}

	playerCount := len(playerLabels)

	// Bucket the neutral zones by tier so low zones always land in the outer band
	// and high zones in the inner band, regardless of the pool's incoming order.
	tierBuckets := [3][]int{}
	for index, plan := range neutralZones {
		switch plan.Quality {
		case neutralZone.QualityHigh:
			tierBuckets[2] = append(tierBuckets[2], index)
		case neutralZone.QualityMedium:
			tierBuckets[1] = append(tierBuckets[1], index)
		case neutralZone.QualityLow:
			fallthrough
		default:
			tierBuckets[0] = append(tierBuckets[0], index)
		}
	}
	// Spread each tier evenly across the players so every fractal grows at the
	// same rate and the overall pattern stays rotationally balanced.
	perPlayerTier := [3][][]int{
		distributeRoundRobin(tierBuckets[0], playerCount),
		distributeRoundRobin(tierBuckets[1], playerCount),
		distributeRoundRobin(tierBuckets[2], playerCount),
	}

	sectorHalf := math.Pi
	if playerCount >= 1 {
		sectorHalf = math.Pi / float64(playerCount)
	}

	var allLabels []string
	var positions models.Positions
	trees := make([]fractalTree, playerCount)

	for player := range playerCount {
		axis := startAngle + 2.0*math.Pi*float64(player)/float64(playerCount)

		// Player zone: the base of this fractal on the outer ring.
		playerIndex := len(allLabels)
		allLabels = append(allLabels, playerLabels[player])
		positions.Add(circlePoint(axis, centreX, centreY, playerRadius))

		var tree fractalTree
		tree.player = playerIndex
		for tier := range 3 {
			plan := perPlayerTier[tier][player]
			half := sectorHalf * tierSpread[tier]
			radius := tierRadius[tier]
			for slot, planIndex := range plan {
				angle := axis
				if len(plan) > 1 {
					angle = axis - half + 2.0*half*float64(slot)/float64(len(plan)-1)
				}
				tree.levels[tier] = append(tree.levels[tier], len(allLabels))
				allLabels = append(allLabels, neutralZones[planIndex].Label)
				positions.Add(circlePoint(angle, centreX, centreY, radius))
			}
		}
		trees[player] = tree
	}

	pairs := this.createFractalPairs(trees)
	return allLabels, positions, pairs
}

// createFractalPairs wires each fractal from the player outward through its tiers
// and then stitches the fractals together at their centre-facing tips.
func (this *FractalTopologyService) createFractalPairs(trees []fractalTree) []models.ConnectionIndexes {
	builder := newPairBuilder()

	tips := make([]int, len(trees))
	for treeIndex, tree := range trees {
		// Collect the tiers that actually received zones, outermost first.
		chain := make([][]int, 0, len(tree.levels))
		for tier := range len(tree.levels) {
			if len(tree.levels[tier]) > 0 {
				chain = append(chain, tree.levels[tier])
			}
		}

		// A player with no neutral zones at all has only itself as a tip; it is
		// left for CreateMissingConnections to attach so no player-player edge is
		// forced here.
		if len(chain) == 0 {
			tips[treeIndex] = tree.player
			continue
		}

		// The player branches into every zone of its outermost tier.
		for _, zoneIndex := range chain[0] {
			builder.add(tree.player, zoneIndex)
		}
		// Each inner tier merges into the tier just outside it: zones stay sorted
		// by angle, so mapping every inner zone to a proportional outer parent
		// produces clean, non-crossing branches that funnel toward the centre.
		for level := 0; level+1 < len(chain); level++ {
			outer := chain[level]
			inner := chain[level+1]
			for innerSlot, innerIndex := range inner {
				parent := outer[innerSlot*len(outer)/len(inner)]
				builder.add(parent, innerIndex)
			}
		}

		// The representative tip is the middle zone of the innermost tier - the
		// point where this fractal reaches toward its neighbours.
		deepest := chain[len(chain)-1]
		tips[treeIndex] = deepest[len(deepest)/2]
	}

	// Weave the tips into a closed ring around the centre. This is the seam where
	// neighbouring fractals integrate, and it also guarantees the whole map is one
	// connected component. Tips that are still a player zone (a fractal with no
	// neutrals) are dropped so the ring never links two players directly.
	neutralTips := make([]int, 0, len(tips))
	for treeIndex, tip := range tips {
		if tip == trees[treeIndex].player {
			continue
		}
		neutralTips = append(neutralTips, tip)
	}
	connectRing(builder, neutralTips)

	return builder.pairs
}

// distributeRoundRobin splits items across the given number of buckets in
// round-robin order, so each bucket differs in size by at most one. Returns a
// slice of empty buckets when there are no buckets to fill.
func distributeRoundRobin(items []int, buckets int) [][]int {
	out := make([][]int, max(0, buckets))
	if buckets <= 0 {
		return out
	}
	for offset, item := range items {
		bucket := offset % buckets
		out[bucket] = append(out[bucket], item)
	}
	return out
}
