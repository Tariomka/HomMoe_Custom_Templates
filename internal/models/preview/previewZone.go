package preview

import "image"

// Zone is one zone laid out on the preview canvas.
type Zone struct {
	Name      string
	Letter    string
	Center    image.Point
	IsPlayer  bool
	IsHub     bool
	Tier      int // 0 unknown, 1 bronze, 2 silver, 3 gold
	Owner     int
	HasCastle bool
	Castles   int
}

type ZoneTier int

const (
	TierUnknown ZoneTier = iota
	TierBronze
	TierSilver
	TierGold
)
