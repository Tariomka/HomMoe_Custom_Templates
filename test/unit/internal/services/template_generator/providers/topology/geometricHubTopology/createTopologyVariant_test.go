package geometricHubTopology_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomPortalsEnabled_AddsExtraPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans([]string{"E", "F", "G", "H"}, nil, nil)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	configuration.PlayerCount = len(playerLabels)
	configuration.RandomPortals = true
	configuration.MaxPortalConnections = 4
	tuning := models.NewGenerationTuning(configuration, len(playerLabels)+len(plans)+1)
	baseline := buildGeoHubVariant(playerLabels, plans)

	// Act
	variant := topology.NewGeometricHubTopologyService().
		CreateTopologyVariant(*configuration, playerLabels, plans, tuning, false)

	// Assert
	assert.Greater(t, len(variant.Connections), len(baseline.Connections))
}

func TestWhenHubMandatoryContentConfigured_HubZoneReferencesHubContentGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans([]string{"E", "F"}, nil, nil)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	configuration.PlayerCount = len(playerLabels)
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "hub_item"}}
	tuning := models.NewGenerationTuning(configuration, len(playerLabels)+len(plans)+1)

	// Act
	variant := topology.NewGeometricHubTopologyService().
		CreateTopologyVariant(*configuration, playerLabels, plans, tuning, false)

	// Assert
	assert.Contains(t, variant.Zones[0].MandatoryContent, "mandatory_content_hub")
}

func TestWhenNoNeutralZonesSelected_EveryPlayerConnectsToHub(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}

	// Act
	variant := buildGeoHubVariant(playerLabels, neutralZone.Plans{})

	// Assert
	assert.ElementsMatch(t,
		[]string{"Spawn-A", "Spawn-B", "Spawn-C", "Spawn-D"},
		hubPortalTargets(variant))
}

func TestWhenNoNeutralZonesSelected_NoPlayerToPlayerConnectionsExist(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}

	// Act
	variant := buildGeoHubVariant(playerLabels, neutralZone.Plans{})

	// Assert
	playerToPlayer := 0
	for _, connection := range variant.Connections {
		if strings.HasPrefix(connection.From, "Spawn-") && strings.HasPrefix(connection.To, "Spawn-") {
			playerToPlayer++
		}
	}
	assert.Zero(t, playerToPlayer)
}

func TestWhenSingleNeutralZoneSelected_SharedZoneConnectsBothAdjacentPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	plans := mixedPlans([]string{"E"}, nil, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	assert.ElementsMatch(t, []string{"Spawn-A", "Spawn-B"}, spawnNeighborsOf(variant, "Neutral-E"))
}

func TestWhenSingleNeutralZoneSelected_UncoveredPlayersConnectDirectlyToHub(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	plans := mixedPlans([]string{"E"}, nil, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	assert.ElementsMatch(t, []string{"Neutral-E", "Spawn-C", "Spawn-D"}, hubPortalTargets(variant))
}

func TestWhenTwoStablesFillEveryGapWithoutCorners_AllStablesConnectToHub(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans([]string{"E", "F", "G", "H"}, nil, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	assert.ElementsMatch(t,
		[]string{"Neutral-E", "Neutral-F", "Neutral-G", "Neutral-H"},
		hubPortalTargets(variant))
}

func TestWhenFullHexagonShapeIsReached_OnlyCornersTouchHub(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans([]string{"E", "F", "G", "H"}, []string{"I", "J"}, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	assert.ElementsMatch(t, []string{"Neutral-I", "Neutral-J"}, hubPortalTargets(variant))
}

func TestWhenFullHexagonShapeIsReached_CornersTakeLowestQualityPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans([]string{"E", "F", "G", "H"}, []string{"I", "J"}, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// The corner zones (hub-adjacent) must be the two Low plans; the Medium
	// stables must sit next to the players instead.
	assert.Empty(t, spawnNeighborsOf(variant, "Neutral-I"),
		"low-quality corner zones must not touch a player zone")
}

func TestWhenOneInteriorZoneExists_InteriorTakesHighestPlanAndSplitsHexagon(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans([]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"K"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Rule 7: the interior (highest plan K) connects the hexagon's two stable
	// zones and the Hub - never the player.
	assert.ElementsMatch(t, []string{"Hub", "Neutral-H", "Neutral-E"}, neighborsOf(variant, "Neutral-K"))
}

func TestWhenHexagonHasTwoInteriors_FormsChainBetweenStables(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"M", "N", "O", "P"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A holds interiors x1=M, x2=O between stables sL=H and sR=E:
	// the chain runs H - M - O - E, and M additionally portals to the hub.
	assert.ElementsMatch(t, []string{"Hub", "Neutral-O", "Neutral-H"}, neighborsOf(variant, "Neutral-M"))
}

func TestWhenHexagonHasTwoInteriors_BothInteriorsPortalToHub(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"M", "N", "O", "P"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	assert.Subset(t, hubPortalTargets(variant),
		[]string{"Neutral-M", "Neutral-N", "Neutral-O", "Neutral-P"})
}

func TestWhenHexagonHasThreeInteriors_PlayerSideInteriorConnectsBothStables(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"M", "N", "O", "P", "Q", "R"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A's triangle is x1=M, x2=O, x3=Q; the player-side vertex x3
	// carries the allowed 4-connection exception: both ring edges plus BOTH
	// stables (sL=H, sR=E) - and no hub portal.
	assert.ElementsMatch(t,
		[]string{"Neutral-M", "Neutral-O", "Neutral-H", "Neutral-E"},
		neighborsOf(variant, "Neutral-Q"))
}

func TestWhenHexagonHasThreeInteriors_PlayerSideInteriorHasNoHubPortal(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"M", "N", "O", "P", "Q", "R"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	assert.NotContains(t, hubPortalTargets(variant), "Neutral-Q")
}

func TestWhenHexagonHasFourInteriors_RingHasNoDiagonals(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"},
		[]string{"M", "N", "O", "P", "Q", "R", "S", "T"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A's square is x1=M, x2=O, x3=Q, x4=S. x1 must connect its ring
	// neighbors x2 and x3 plus sL=H and the hub - never the diagonal x4=S.
	assert.ElementsMatch(t,
		[]string{"Hub", "Neutral-O", "Neutral-Q", "Neutral-H"},
		neighborsOf(variant, "Neutral-M"))
}

func TestWhenHexagonHasFourInteriors_PlayerSideChainLinksStables(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"},
		[]string{"M", "N", "O", "P", "Q", "R", "S", "T"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// The player-side pair x3=Q, x4=S forms the chain sL - x3 - x4 - sR:
	// x3 connects x1=M (ring), x4=S (ring) and sL=H - no hub portal.
	assert.ElementsMatch(t,
		[]string{"Neutral-M", "Neutral-S", "Neutral-H"},
		neighborsOf(variant, "Neutral-Q"))
}

func TestWhenHexagonHasFiveInteriors_RingConnectsAngularNeighborsOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"},
		[]string{"M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A's pentagon is x1=M, x2=O, x3=Q, x4=S, x5=U. The player-axis
	// vertex x5 touches only its ring neighbors x3 and x4 - no stable, no hub.
	assert.ElementsMatch(t, []string{"Neutral-Q", "Neutral-S"}, neighborsOf(variant, "Neutral-U"))
}

func TestWhenHexagonHasFiveInteriors_StablesLinkNearestSideVertices(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"},
		[]string{"M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// sL=H connects the hub-side-left vertex x1=M plus its nearest left-side
	// vertex x3=Q; x3 in turn touches its ring neighbors x1=M and x5=U.
	assert.ElementsMatch(t,
		[]string{"Neutral-M", "Neutral-U", "Neutral-H"},
		neighborsOf(variant, "Neutral-Q"))
}

func TestWhenHexagonHasFiveInteriors_NoZoneOutsideHubExceedsFourConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"},
		[]string{"M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	overConnected := map[string]int{}
	for _, zone := range variant.Zones {
		if zone.Name == "Hub" {
			continue
		}
		if count := len(neighborsOf(variant, zone.Name)); count > 4 {
			overConnected[zone.Name] = count
		}
	}
	assert.Empty(t, overConnected, "only the Hub may exceed four connections")
}

func TestWhenInteriorsAssigned_HubFacingVerticesTakeHighestPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"M", "N", "O", "P", "Q", "R"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hub portals are the merged corners (lowest plans I, J) plus the x1/x2
	// vertices of both hexagons - the four highest plans M, N, O, P. The two
	// remaining highs Q, R sit on the player-side vertices without a portal.
	assert.ElementsMatch(t,
		[]string{"Neutral-I", "Neutral-J", "Neutral-M", "Neutral-N", "Neutral-O", "Neutral-P"},
		hubPortalTargets(variant))
}

func TestWhenThreePlayersHaveFifteenNeutrals_HubPortalsAreCornersAndHubFacingInteriors(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C"}
	plans := mixedPlans(
		[]string{"D", "E", "F", "G", "H", "I"}, []string{"J", "K", "L"},
		[]string{"M", "N", "O", "P", "Q", "R"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// 15 neutrals over 3 players: 6 stables, 3 merged corners, 6 interiors
	// (a k=2 chain per hexagon). The hub sees the 3 corners plus all 6
	// interiors - the merged corners are kept, no splits exist.
	assert.ElementsMatch(t,
		[]string{
			"Neutral-J", "Neutral-K", "Neutral-L",
			"Neutral-M", "Neutral-N", "Neutral-O", "Neutral-P", "Neutral-Q", "Neutral-R",
		},
		hubPortalTargets(variant))
}

func TestWhenAnyConnectionTouchesHub_ItIsPortalType(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C"}
	plans := mixedPlans(
		[]string{"D", "E", "F", "G", "H", "I"}, []string{"J", "K"}, []string{"L", "M"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	nonPortal := 0
	for _, connection := range hubConnections(variant) {
		if connection.ConnectionType != "Portal" {
			nonPortal++
		}
	}
	assert.Zero(t, nonPortal, "every connection touching the Hub must be a portal")
}

func TestWhenVariantIsBuilt_EveryZoneCarriesGeneratorPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C"}
	plans := mixedPlans(
		[]string{"D", "E", "F", "G", "H", "I"}, []string{"J", "K"}, []string{"L", "M"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	missing := 0
	for _, zone := range variant.Zones {
		if zone.GeneratorPosition == nil {
			missing++
		}
	}
	assert.Zero(t, missing, "every zone must be stamped with a fixed-geometry position")
}

func TestWhenVariantIsBuilt_CreatesOneZonePerPlanPlusPlayersPlusHub(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C"}
	plans := mixedPlans(
		[]string{"D", "E", "F", "G", "H", "I"}, []string{"J", "K"}, []string{"L", "M"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	assert.Len(t, variant.Zones, len(playerLabels)+len(plans)+1)
}

func TestWhenEveryConnectionIsInspected_AllEndpointsReferenceExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H", "I", "J", "K", "L"}, []string{"M", "N", "O", "P"}, []string{"Q", "R"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	zoneNames := map[string]bool{}
	for _, zone := range variant.Zones {
		zoneNames[zone.Name] = true
	}
	dangling := 0
	for _, connection := range variant.Connections {
		if !zoneNames[connection.From] || !zoneNames[connection.To] {
			dangling++
		}
	}
	assert.Zero(t, dangling)
}

func TestWhenThreePlayersFillFullHexagons_HexagonSidesAreEqualLength(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C"}
	plans := mixedPlans([]string{"D", "E", "F", "G", "H", "I"}, []string{"J", "K", "L"}, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A's perimeter is Hub - cL(J) - sL(I) - player - sR(D) - cR(L) -
	// Hub; a regular hexagon means all six sides share one length.
	sides := []float64{
		distanceBetween(variant, "Hub", "Neutral-J"),
		distanceBetween(variant, "Neutral-J", "Neutral-I"),
		distanceBetween(variant, "Neutral-I", "Spawn-A"),
		distanceBetween(variant, "Spawn-A", "Neutral-D"),
		distanceBetween(variant, "Neutral-D", "Neutral-L"),
		distanceBetween(variant, "Neutral-L", "Hub"),
	}
	assert.InDelta(t, 0, spreadOf(sides), 0.000001,
		"3-player hexagons must be regular: %v", sides)
}

func TestWhenFourPlayersFillFullHexagons_UnsharedVerticesAreExactly120Degrees(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H", "I", "J", "K", "L"}, []string{"M", "N", "O", "P"}, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A's perimeter is Hub - cL(M) - sL(L) - player - sR(E) - cR(P).
	// The hub angle is locked to the 90-degree sector, so the stables and the
	// player keep the regular hexagon's 120 degrees while the two shared
	// corners absorb the forced surplus (135 degrees each).
	angles := perimeterFreeAngles(variant,
		[6]string{"Hub", "Neutral-M", "Neutral-L", "Spawn-A", "Neutral-E", "Neutral-P"})
	assert.InDeltaSlice(t, []float64{135, 120, 120, 120, 135}, angles, 0.000001,
		"4-player hexagons must keep 120 degrees on every unshared vertex: %v", angles)
}

func TestWhenFivePlayersFillFullHexagons_UnsharedVerticesAreExactly120Degrees(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D", "E"}
	plans := mixedPlans(
		[]string{"F", "G", "H", "I", "J", "K", "L", "M", "N", "O"},
		[]string{"P", "Q", "R", "S", "T"}, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A's perimeter is Hub - cL(P) - sL(O) - player - sR(F) - cR(T).
	// The hub angle is locked to the 72-degree sector, so the stables and the
	// player keep the regular hexagon's 120 degrees while the two shared
	// corners absorb the forced surplus (144 degrees each).
	angles := perimeterFreeAngles(variant,
		[6]string{"Hub", "Neutral-P", "Neutral-O", "Spawn-A", "Neutral-F", "Neutral-T"})
	assert.InDeltaSlice(t, []float64{144, 120, 120, 120, 144}, angles, 0.000001,
		"5-player hexagons must keep 120 degrees on every unshared vertex: %v", angles)
}

func TestWhenFewPlayersCompareToMany_FewPlayersSitCloserToHub(t *testing.T) {
	t.Parallel()
	// Arrange
	threePlayerLabels := []string{"A", "B", "C"}
	sixPlayerLabels := []string{"A", "B", "C", "D", "E", "F"}

	// Act
	threePlayerVariant := buildGeoHubVariant(threePlayerLabels, neutralZone.Plans{})
	sixPlayerVariant := buildGeoHubVariant(sixPlayerLabels, neutralZone.Plans{})

	// Assert
	assert.Less(t,
		distanceBetween(threePlayerVariant, "Hub", "Spawn-A"),
		distanceBetween(sixPlayerVariant, "Hub", "Spawn-A"))
}

func TestWhenHexagonHasThreeInteriors_VerticesAreEquidistantFromPolygonCenter(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"M", "N", "O", "P", "Q", "R"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// Hexagon A's triangle vertices M, O, Q must lie on one circle around
	// their centroid (the regular-polygon center).
	vertexM := positionOf(variant, "Neutral-M")
	vertexO := positionOf(variant, "Neutral-O")
	vertexQ := positionOf(variant, "Neutral-Q")
	centroid := [2]float64{
		(vertexM[0] + vertexO[0] + vertexQ[0]) / 3,
		(vertexM[1] + vertexO[1] + vertexQ[1]) / 3,
	}
	radii := make([]float64, 0, 3)
	for _, vertex := range [][2]float64{vertexM, vertexO, vertexQ} {
		radii = append(radii, math.Hypot(vertex[0]-centroid[0], vertex[1]-centroid[1]))
	}
	assert.InDelta(t, 0, spreadOf(radii), 0.000001,
		"triangle vertices must be equidistant from the polygon center: %v", radii)
}

func TestWhenHexagonHasTwoInteriors_ChainSpacingIsEven(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H"}, []string{"I", "J"}, []string{"M", "N", "O", "P"})

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// The hexagon-A chain sL(H) - x1(M) - x2(O) - sR(E) must be evenly spaced.
	spacings := []float64{
		distanceBetween(variant, "Neutral-H", "Neutral-M"),
		distanceBetween(variant, "Neutral-M", "Neutral-O"),
		distanceBetween(variant, "Neutral-O", "Neutral-E"),
	}
	assert.InDelta(t, 0, spreadOf(spacings), 0.001,
		"k=2 chain must be evenly spaced: %v", spacings)
}

func TestWhenFourPlayerHexagonHasTwoInteriors_ChainSpacingIsEven(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	plans := mixedPlans(
		[]string{"E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
			"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X"}, nil, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	// The equiangular geometry solves the k=2 circumradius from the actual
	// stable positions: hexagon A's chain sL(T) - x1(E) - x2(I) - sR(M) must
	// stay evenly spaced under the 4-player blueprint too.
	spacings := []float64{
		distanceBetween(variant, "Neutral-T", "Neutral-E"),
		distanceBetween(variant, "Neutral-E", "Neutral-I"),
		distanceBetween(variant, "Neutral-I", "Neutral-M"),
	}
	assert.InDelta(t, 0, spreadOf(spacings), 0.000001,
		"4-player k=2 chain must be evenly spaced: %v", spacings)
}

func TestWhenEightPlayersHaveManyInteriors_AllPositionsStayInsideUnitSquare(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	neutralLabels := make([]string, 64)
	for index := range neutralLabels {
		neutralLabels[index] = fmt.Sprintf("N%02d", index)
	}
	plans := mixedPlans(neutralLabels, nil, nil)

	// Act
	variant := buildGeoHubVariant(playerLabels, plans)

	// Assert
	outOfBounds := 0
	for _, zone := range variant.Zones {
		position := *zone.GeneratorPosition
		if position[0] < 0 || position[0] > 1 || position[1] < 0 || position[1] > 1 {
			outOfBounds++
		}
	}
	assert.Zero(t, outOfBounds, "every position must stay inside the normalized unit square")
}
