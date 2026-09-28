package connectionCurveLayout_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenASingleConnectionSpansAClearChord_ItsControlPointStaysOnTheMidpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordPositions())

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 350.0)}, controlPoints(curves))
}

func TestWhenTwoUnobstructedConnectionsSharePair_TheyFanOutTwentyOnePixelsApart(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordPositions())

	// Act
	curves := layout.Build(parallelConnections(2))

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{data.NewVec2(350.0, 371.0), data.NewVec2(350.0, 329.0)},
		controlPoints(curves))
}

func TestWhenThreeUnobstructedConnectionsSharePair_TheMiddleOneStaysOnTheChord(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordPositions())

	// Act
	curves := layout.Build(parallelConnections(3))

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{
			data.NewVec2(350.0, 392.0),
			data.NewVec2(350.0, 350.0),
			data.NewVec2(350.0, 308.0),
		},
		controlPoints(curves))
}

func TestWhenAConnectionIsReversed_ItFansWithItsForwardTwinAroundTheSameChord(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordPositions())
	connections := []template_model.Connection{newConnection("A", "B"), newConnection("B", "A")}

	// Act
	curves := layout.Build(connections)

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{data.NewVec2(350.0, 371.0), data.NewVec2(350.0, 329.0)},
		controlPoints(curves))
}

func TestWhenAConnectionIsReversed_ItsCurveRunsInCanonicalEndpointOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordPositions())

	// Act
	curves := layout.Build([]template_model.Connection{newConnection("B", "A")})

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{data.NewVec2(140.0, 350.0), data.NewVec2(560.0, 350.0)},
		[]data.Vec2[float64]{curves[0].Start, curves[0].End})
}

func TestWhenConnectionsSharePairs_TheCurvesAreGroupedInFirstSeenOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := chordPositions()
	positions["C"] = data.NewVec2(350.0, 140.0)
	layout := newLayout(positions)
	connections := []template_model.Connection{
		newConnection("A", "B"),
		newConnection("A", "C"),
		newConnection("B", "A"),
	}

	// Act
	curves := layout.Build(connections)

	// Assert
	assert.Equal(t, []int{0, 2, 1}, connectionIndices(curves))
}

func TestWhenAnEndpointHasNoPosition_TheConnectionIsSkipped(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordPositions())
	connections := []template_model.Connection{newConnection("A", "B"), newConnection("A", "Z")}

	// Act
	curves := layout.Build(connections)

	// Assert
	assert.Equal(t, []int{0}, connectionIndices(curves))
}

func TestWhenBothEndpointsShareAPosition_TheCurveCollapsesOntoThatPoint(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(map[string]data.Vec2[float64]{
		"A": data.NewVec2(140.0, 350.0),
		"B": data.NewVec2(140.0, 350.0),
	})

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(140.0, 350.0)}, controlPoints(curves))
}

func TestWhenTheChordIsShorterThanAPixel_NoObstacleBendsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(map[string]data.Vec2[float64]{
		"A": data.NewVec2(140.0, 350.0),
		"B": data.NewVec2(140.5, 350.0),
		"D": data.NewVec2(140.25, 350.0),
	})

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(140.25, 350.0)}, controlPoints(curves))
}

func TestWhenAZoneSitsJustBelowTheChord_TheCurveTakesTheShorterWayAboveIt(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 364.0))

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 274.0)}, controlPoints(curves))
}

func TestWhenAZoneSitsJustAboveTheChord_TheCurveTakesTheShorterWayBelowIt(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 336.0))

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 426.0)}, controlPoints(curves))
}

func TestWhenAZoneSitsExactlyOnTheChord_TheTieBendsTheCurveToThePositiveSide(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 350.0))

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 246.0)}, controlPoints(curves))
}

// The review's reproduction: two equal obstacles mirrored across the chord used
// to resolve by map iteration order, yielding two different curves for one input.
func TestWhenEqualObstaclesMirrorEachOther_EveryBuildBendsToThePositiveSide(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{newConnection("A", "B")}
	seen := map[data.Vec2[float64]]bool{}

	// Act
	for attempt := range 200 {
		positions := map[string]data.Vec2[float64]{}
		names := []string{"A", "B", "C", "D"}
		points := []data.Vec2[float64]{
			data.NewVec2(100.0, 350.0),
			data.NewVec2(600.0, 350.0),
			data.NewVec2(350.0, 340.0),
			data.NewVec2(350.0, 360.0),
		}
		for offset := range names {
			index := (attempt + offset) % len(names)
			positions[names[index]] = points[index]
		}
		layout := preview.ConnectionCurveLayout{Positions: positions, ZoneRadius: 21.0}
		seen[layout.Build(connections)[0].Control] = true
	}

	// Assert
	assert.Equal(t, map[data.Vec2[float64]]bool{data.NewVec2(350.0, 260.0): true}, seen)
}

func TestWhenObstaclesSitOnBothSides_TheCurveTakesTheShorterDetourThatClearsBoth(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := chordWithObstacle(350.0, 340.0)
	positions["E"] = data.NewVec2(350.0, 370.0)
	layout := newLayout(positions)

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 226.0)}, controlPoints(curves))
}

func TestWhenTwoConnectionsAreObstructed_OnePassesOnEachSide(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 364.0))

	// Act
	curves := layout.Build(parallelConnections(2))

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{data.NewVec2(350.0, 482.0), data.NewVec2(350.0, 274.0)},
		controlPoints(curves))
}

func TestWhenThreeConnectionsAreObstructedAndThePositiveSideIsShorter_TwoPassAbove(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 364.0))

	// Act
	curves := layout.Build(parallelConnections(3))

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{
			data.NewVec2(350.0, 482.0),
			data.NewVec2(350.0, 274.0),
			data.NewVec2(350.0, 232.0),
		},
		controlPoints(curves))
}

func TestWhenThreeConnectionsAreObstructedAndTheNegativeSideIsShorter_TwoPassBelow(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 336.0))

	// Act
	curves := layout.Build(parallelConnections(3))

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{
			data.NewVec2(350.0, 468.0),
			data.NewVec2(350.0, 426.0),
			data.NewVec2(350.0, 218.0),
		},
		controlPoints(curves))
}

func TestWhenThreeConnectionsAreObstructedAndBothSidesTie_TwoPassOnThePositiveSide(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 350.0))

	// Act
	curves := layout.Build(parallelConnections(3))

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{
			data.NewVec2(350.0, 454.0),
			data.NewVec2(350.0, 246.0),
			data.NewVec2(350.0, 204.0),
		},
		controlPoints(curves))
}

func TestWhenFourConnectionsAreObstructed_EachSideStepsOutwards(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 364.0))

	// Act
	curves := layout.Build(parallelConnections(4))

	// Assert
	assert.Equal(t,
		[]data.Vec2[float64]{
			data.NewVec2(350.0, 524.0),
			data.NewVec2(350.0, 482.0),
			data.NewVec2(350.0, 274.0),
			data.NewVec2(350.0, 232.0),
		},
		controlPoints(curves))
}

// The detour is a midpoint offset, so an obstacle a fifth of the way along the
// chord pushes the curve exactly as far as one at its centre would.
func TestWhenAnObstacleSitsOffCentre_TheDetourIsAppliedAtTheMidpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(224.0, 364.0))

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 274.0)}, controlPoints(curves))
}

func TestWhenAnObstacleSitsInsideTheChordMargin_ItIsIgnored(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(161.0, 364.0))

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 350.0)}, controlPoints(curves))
}

func TestWhenAnObstacleSitsAtTheClearance_ItIsIgnored(t *testing.T) {
	t.Parallel()
	// Arrange
	layout := newLayout(chordWithObstacle(350.0, 396.0))

	// Act
	curves := layout.Build(parallelConnections(1))

	// Assert
	assert.Equal(t, []data.Vec2[float64]{data.NewVec2(350.0, 350.0)}, controlPoints(curves))
}
