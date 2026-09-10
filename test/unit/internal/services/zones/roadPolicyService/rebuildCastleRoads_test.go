package roadPolicyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenZoneGainedCastles_AddsTheMissingPrimaryRoutes(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 3)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []string{"1", "2"}, roadTargets(zone, "MainObject"))
}

func TestWhenZoneLostCastles_DropsTheRoutesToThem(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 1,
		template_model.Road{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("1")},
		template_model.Road{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("2")})

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Empty(t, zone.Roads)
}

func TestWhenPrimaryRouteAlreadyExists_DoesNotDuplicateIt(t *testing.T) {
	t.Parallel()
	// Arrange
	existing := template_model.Road{
		Type:       "Stone",
		From:       mainObjectRef("0"),
		To:         mainObjectRef("1"),
		GuardValue: 1500,
	}
	zone := castleZone("A", 2, existing)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []template_model.Road{existing}, zone.Roads)
}

func TestWhenPrimaryRouteExistsInReverseOrder_DoesNotDuplicateIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 2,
		template_model.Road{Type: "Stone", From: mainObjectRef("1"), To: mainObjectRef("0")})

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Len(t, zone.Roads, 1)
}

func TestWhenZoneHasACustomRouteBetweenOtherCastles_KeepsItAndStillAddsThePrimaryOnes(t *testing.T) {
	t.Parallel()
	// Arrange
	custom := template_model.Road{Type: "Dirt", From: mainObjectRef("1"), To: mainObjectRef("2")}
	zone := castleZone("A", 3, custom)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []template_model.Road{
		custom,
		{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("1")},
		{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("2")},
	}, zone.Roads)
}

func TestWhenZoneHasNonCastleRoads_LeavesThemUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	approach := road(mainObjectRef("0"), connectionRef("Rnd-A-B"))
	foothold := road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))
	zone := castleZone("A", 1, approach, foothold)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []template_model.Road{approach, foothold}, zone.Roads)
}

func TestWhenCastleRouteAnchorIsNotAnIndex_KeepsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	opaque := template_model.Road{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("primary")}
	zone := castleZone("A", 1, opaque)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []template_model.Road{opaque}, zone.Roads)
}

func TestWhenZoneHasNoMainObjects_CreatesNoCastleRoute(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 0)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Empty(t, zone.Roads)
}

func TestWhenNothingChangesAndTheRoadListIsNil_LeavesItNil(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 1)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Nil(t, zone.Roads)
}

func TestWhenNothingChangesAndTheRoadListIsEmpty_LeavesItEmptyNotNil(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 1)
	zone.Roads = []template_model.Road{}

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.NotNil(t, zone.Roads)
}

// ── Arena marker ─────────────────────────────────────────────────────

func TestWhenTheArenaMarkerFollowsTheCastles_CreatesNoRouteToTheMarker(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 2)
	zone.MainObjects = append(zone.MainObjects, arenaMarker())

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []template_model.Road{
		{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("1")},
	}, zone.Roads)
}

// An imported template may list the marker first, in which case the primary
// route runs between the real castles behind it.
func TestWhenTheArenaMarkerPrecedesTheCastles_RoutesBetweenTheCastlesBehindIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 2)
	zone.MainObjects = append([]template_model.MainObject{arenaMarker()}, zone.MainObjects...)

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []template_model.Road{
		{Type: "Stone", From: mainObjectRef("1"), To: mainObjectRef("2")},
	}, zone.Roads)
}

func TestWhenZoneHoldsOnlyTheArenaMarker_CreatesNoCastleRoute(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := arenaZone()

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Empty(t, zone.Roads)
}

func TestWhenTheArenaMarkerSitsBetweenTheCastles_RoutesAroundIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := castleZone("A", 2)
	zone.MainObjects = []template_model.MainObject{
		zone.MainObjects[0],
		arenaMarker(),
		zone.MainObjects[1],
	}

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, []template_model.Road{
		{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("2")},
	}, zone.Roads)
}

// ── Caller aliasing ──────────────────────────────────────────────────

// The dialog hands over roads it shares with the previewed base, so dropping a
// stale route must not compact the caller's array underneath it.
func TestWhenAStaleRouteIsDropped_LeavesTheCallersArrayIntact(t *testing.T) {
	t.Parallel()
	// Arrange
	stale := template_model.Road{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("2")}
	shared := []template_model.Road{
		{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("1")},
		stale,
	}
	zone := castleZone("A", 2)
	zone.Roads = shared

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, stale, shared[1])
}

// Appending a missing route must allocate rather than claim whatever spare
// capacity the caller's slice still had.
func TestWhenARouteIsAppended_LeavesTheCallersSpareCapacityIntact(t *testing.T) {
	t.Parallel()
	// Arrange
	spare := road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))
	shared := make([]template_model.Road, 2, 4)
	shared[0] = road(mainObjectRef("0"), connectionRef("Rnd-A-B"))
	shared[1] = spare
	zone := castleZone("A", 2)
	zone.Roads = shared[:1]

	// Act
	newPolicy().RebuildCastleRoads(&zone)

	// Assert
	assert.Equal(t, spare, shared[1])
}

// ── Unresolvable factory anchors ─────────────────────────────────────

// The real factory numbers its anchors 0..n-1, so a candidate whose anchor
// cannot be resolved only ever arrives from another IRoadFactory. The marker
// in front of the castle is what forces the rebase to run at all.
func TestWhenAFactoryCandidateAnchorCannotBeResolved_KeepsTheReferenceAsItCame(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		anchor template_model.TypedRef
	}{
		{name: "WhenTheAnchorIsNotANumber_KeepsTheReferenceAsItCame", anchor: mainObjectRef("primary")},
		{name: "WhenTheAnchorIsNegative_KeepsTheReferenceAsItCame", anchor: mainObjectRef("-1")},
		{name: "WhenTheAnchorIsPastTheLastAnchor_KeepsTheReferenceAsItCame", anchor: mainObjectRef("7")},
		{
			name:   "WhenTheAnchorCarriesNoArguments_KeepsTheReferenceAsItCame",
			anchor: template_model.TypedRef{Type: "MainObject"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			candidate := template_model.Road{
				Type: "Stone",
				From: testCase.anchor,
				To:   connectionRef("Rnd-A-B"),
			}
			factory := &test_helpers.RoadFactoryMock{}
			factory.On("CreateOuterZoneRoads", mock.Anything, 1, 0, false).
				Return([]template_model.Road{candidate})
			zone := castleZone("A", 1)
			zone.MainObjects = append([]template_model.MainObject{arenaMarker()}, zone.MainObjects...)

			// Act
			zones.NewRoadPolicyService(factory).RebuildCastleRoads(&zone)

			// Assert
			assert.Equal(t, []template_model.Road{candidate}, zone.Roads)
		})
	}
}
