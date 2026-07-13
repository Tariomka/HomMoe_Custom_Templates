package common

type MapSize struct {
	Size  int
	Label string
}

var BaseMapSizes = []MapSize{
	{Size: 64, Label: "64x64 - S"},
	{Size: 80, Label: "80x80 - M"},
	{Size: 96, Label: "96x96 - M"},
	{Size: 112, Label: "112x112 - L"},
	{Size: 128, Label: "128x128 - L"},
	{Size: 144, Label: "144x144 - XL"},
	{Size: 160, Label: "160x160 - XL"},
	{Size: 176, Label: "176x176 - H"},
	{Size: 192, Label: "192x192 - H"},
	{Size: 208, Label: "208x208 - G"},
	{Size: 240, Label: "240x240 - G"},
}

var ExpandedMapSizes = []MapSize{
	{Size: 256, Label: "256x256 - C"},
	{Size: 272, Label: "272x272 - C"},
	{Size: 288, Label: "288x288 - C"},
	{Size: 304, Label: "304x304 - C"},
	{Size: 320, Label: "320x320 - C"},
	{Size: 336, Label: "336x336 - C"},
	{Size: 352, Label: "352x352 - C"},
	{Size: 368, Label: "368x368 - C"},
	{Size: 384, Label: "384x384 - C"},
	{Size: 400, Label: "400x400 - C"},
	{Size: 416, Label: "416x416 - C"},
	{Size: 432, Label: "432x432 - C"},
	{Size: 448, Label: "448x448 - C"},
	{Size: 464, Label: "464x464 - C"},
	{Size: 480, Label: "480x480 - C"},
	{Size: 496, Label: "496x496 - C"},
	{Size: 512, Label: "512x512 - C"},
}

var AllMapSizes = append(BaseMapSizes, ExpandedMapSizes...)

func GetMapSize(size int) MapSize {
	for _, mapSize := range AllMapSizes {
		if mapSize.Size == size {
			return mapSize
		}
	}
	return BaseMapSizes[0]
}

func GetMapSizes(withExperimental bool) []MapSize {
	if withExperimental {
		return AllMapSizes
	}
	return BaseMapSizes
}

// GetNearestMapSize returns the valid map size closest to the given size.
// Ties resolve to the smaller size.
func GetNearestMapSize(size int) MapSize {
	nearest := AllMapSizes[0]
	for _, candidate := range AllMapSizes[1:] {
		if absoluteDifference(candidate.Size, size) < absoluteDifference(nearest.Size, size) {
			nearest = candidate
		}
	}
	return nearest
}

func absoluteDifference(left, right int) int {
	if left > right {
		return left - right
	}
	return right - left
}
