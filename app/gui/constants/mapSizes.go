package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type MapSize = models.MapSize

func GetMapSize(size int) MapSize {
	return common.GetMapSize(size)
}

func GetMapSizes(withExperimental bool) []MapSize {
	return common.GetMapSizes(withExperimental)
}
