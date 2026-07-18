package topology

import "math"

// Geometric Hub geometry tuning. Players sit closer to the hub for 2-4
// players (larger hexagons read better) and further out for 5+ (narrow
// sectors need the spread).
const (
	geoHubPlayerRadius          = 0.46
	geoHubPlayerRadiusClose     = 0.38
	geoHubClosePlayerMaximum    = 4
	geoHubEquiangularMinPlayers = 3
	geoHubEquiangularMaxPlayers = 5
	geoHubStableAngleFraction   = 0.35
	geoHubStableIdealOffset     = math.Pi / 6
	// geoHubInteriorPolygonFactor is the interior k-gon circumradius as a
	// fraction of the hexagon side for the legacy layout, chosen so the k=2
	// chain sL-x1-x2-sR is evenly spaced (exact solution ~= 0.3568 for the
	// regular hexagon).
	geoHubInteriorPolygonFactor = 0.357
)

// geoHubGeometry is the polar blueprint of one player hexagon: the radii of
// the player/stable/corner vertices, the stable's angular offset off its
// flanking player's axis, and the interior k-gon center radius/circumradius.
// Corners always sit on the gap mid-angle rays (they are shared between
// adjacent hexagons).
type geoHubGeometry struct {
	playerRadius         float64
	cornerRadius         float64
	stableRadius         float64
	stableOffset         float64
	interiorCenter       float64
	interiorCircumradius float64
}

// newGeoHubGeometry derives the hexagon blueprint for a player count. The P
// hexagon tips meeting at the hub must share the full turn, locking the hub
// angle to the sector (360/P degrees), and the two corners are locked to the
// gap mid-angle rays - so a fully equiangular hexagon only exists for 3
// players. For 3-5 players the closest shape places every UNSHARED vertex
// (both stables and the player) at exactly 120 degrees via the regular-
// hexagon-cap construction (see newHexagonCapGeometry); the two shared
// corners absorb the forced surplus symmetrically (135 degrees for 4P, 144
// for 5P) and the construction reduces to the exact regular hexagon for 3P.
// For 2 players the hub angle (180 degrees) degenerates the construction and
// for 6+ players the pinched sectors would stretch it into slivers, so both
// keep the legacy regular-hexagon-ratio layout (sqrt(3)/2 stable radius,
// +-30 degree ideal offset with a sector-fraction fallback).
func newGeoHubGeometry(playerCount int) geoHubGeometry {
	playerCount = max(playerCount, 2)
	sector := 2 * math.Pi / float64(playerCount)
	playerRadius := geoHubPlayerRadiusFor(playerCount)
	if playerCount >= geoHubEquiangularMinPlayers && playerCount <= geoHubEquiangularMaxPlayers {
		return newHexagonCapGeometry(sector, playerRadius)
	}
	side := playerRadius / 2
	return geoHubGeometry{
		playerRadius:         playerRadius,
		cornerRadius:         side,
		stableRadius:         playerRadius * math.Sqrt(3) / 2,
		stableOffset:         math.Min(geoHubStableIdealOffset, geoHubStableAngleFraction*sector),
		interiorCenter:       side,
		interiorCircumradius: geoHubInteriorPolygonFactor * side,
	}
}

// newHexagonCapGeometry places the five visible zones (cL, sL, player, sR,
// cR) on five consecutive vertices of a TRUE regular hexagon whose sixth,
// omitted vertex faces the hub. Every angle between hexagon edges is then
// the regular 120 degrees - the stables and the player are exactly
// equiangular - and only the hub-corner closing edges deviate, parking the
// geometrically forced surplus on the shared corners. With the regular
// hexagon (side t, center at distance d along the player axis) the corner
// vertex t*(d/t - 1/2, sqrt(3)/2) must lie on the gap mid-ray at halfGap off
// the axis, giving d = t*(1/2 + (sqrt(3)/2)/tan(halfGap)) and the player at
// d + t = playerRadius. For 3 players the omitted vertex IS the hub and the
// full regular hexagon returns.
func newHexagonCapGeometry(sector, playerRadius float64) geoHubGeometry {
	halfGap := sector / 2
	spread := math.Sqrt(3) / 2 / math.Tan(halfGap)
	hexagonSide := playerRadius / (1.5 + spread)
	hexagonCenter := hexagonSide * (0.5 + spread)

	cornerX := hexagonCenter - hexagonSide/2
	stableX := hexagonCenter + hexagonSide/2
	vertexY := hexagonSide * math.Sqrt(3) / 2
	return geoHubGeometry{
		playerRadius:         playerRadius,
		cornerRadius:         math.Hypot(cornerX, vertexY),
		stableRadius:         math.Hypot(stableX, vertexY),
		stableOffset:         math.Atan2(vertexY, stableX),
		interiorCenter:       hexagonCenter,
		interiorCircumradius: evenChainCircumradius(stableX, vertexY, hexagonCenter),
	}
}

// geoHubPlayerRadiusFor picks the player circle radius: 2-4 players sit close
// to the hub, 5+ keep the wider spread that suits narrow sectors.
func geoHubPlayerRadiusFor(playerCount int) float64 {
	if playerCount <= geoHubClosePlayerMaximum {
		return geoHubPlayerRadiusClose
	}
	return geoHubPlayerRadius
}

// evenChainCircumradius solves the even k=2 chain spacing
// sL-x1 == x1-x2 == x2-sR for the actual stable position (stableX, stableY
// in the player-axis frame) and interior polygon center: with x1, x2 at
// (centerX, -+r) the spacing equation 4r^2 == (stableX-centerX)^2 +
// (stableY-r)^2 reduces to the quadratic 3r^2 + 2*stableY*r -
// ((stableX-centerX)^2 + stableY^2) == 0 whose positive root is returned.
func evenChainCircumradius(stableX, stableY, centerX float64) float64 {
	deltaX := stableX - centerX
	return (-stableY + math.Sqrt(4*stableY*stableY+3*deltaX*deltaX)) / 3
}
