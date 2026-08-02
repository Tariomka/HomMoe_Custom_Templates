package common

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

func GetMapSizes(withExperimentalSizes bool) []models.MapSize {
	baseSizes := getBaseMapSizes()
	if withExperimentalSizes {
		return append(baseSizes, getExpandedMapSizes()...)
	}

	return baseSizes
}

func GetMapSize(size int) models.MapSize {
	allSizes := GetMapSizes(true)
	for _, mapSize := range allSizes {
		if mapSize.Size == size {
			return mapSize
		}
	}

	return allSizes[0]
}

func GetNearestMapSize(size int) models.MapSize {
	allSizes := GetMapSizes(true)
	nearest := allSizes[0]
	for _, candidate := range allSizes[1:] {
		if math.Abs(float64(candidate.Size-size)) < math.Abs(float64(nearest.Size-size)) {
			nearest = candidate
		}
	}
	return nearest
}

func getBaseMapSizes() []models.MapSize {
	return []models.MapSize{
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
}

func getExpandedMapSizes() []models.MapSize {
	return []models.MapSize{
		{Size: 256, Label: "256x256 - C"},
		{Size: 272, Label: "272x272 - C"},
		{Size: 288, Label: "288x288 - C"},
		{Size: 304, Label: "304x304 - C"},
		{Size: 320, Label: "320x320 - C"},
		{Size: 336, Label: "336x336 - C"},
		{Size: 352, Label: "352x352 - C"},
		{Size: 368, Label: "368x368 - C"},
		{Size: 384, Label: "384x384 - XC"},
		{Size: 400, Label: "400x400 - XC"},
		{Size: 416, Label: "416x416 - XC"},
		{Size: 432, Label: "432x432 - XC"},
		{Size: 448, Label: "448x448 - XC"},
		{Size: 464, Label: "464x464 - XC"},
		{Size: 480, Label: "480x480 - XC"},
		{Size: 496, Label: "496x496 - XC"},
		{Size: 512, Label: "512x512 - XC"},
	}
}
