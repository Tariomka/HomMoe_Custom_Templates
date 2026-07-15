package preview

import "image"

// Zone is one zone laid out on the preview canvas.
type Zone struct {
	Name      string
	Letter    string
	Center    image.Point
	IsPlayer  bool
	IsHub     bool
	Tier      ZoneTier
	Owner     int
	HasCastle bool
	Castles   int
}

// ZoneTier is the preview quality tier of a zone, ordered lowest to highest.
type ZoneTier int

const (
	TierPlastic ZoneTier = iota
	TierBronze
	TierSilver
	TierGold
	TierPlatinum
)
