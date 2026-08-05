package validators

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type rangedFloatField struct {
	field   floatField
	lowest  float64
	highest float64
}

func newRangedFloatField(
	name string, lowest, highest float64, value func(state *dtos.EditorStateDto) *float64,
) rangedFloatField {
	return rangedFloatField{floatField{name, value}, lowest, highest}
}
