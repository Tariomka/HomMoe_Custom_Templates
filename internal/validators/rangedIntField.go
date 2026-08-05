package validators

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type rangedIntField struct {
	field   intField
	lowest  int
	highest int
}

func newRangedIntField(
	name string, lowest, highest int, value func(state *dtos.EditorStateDto) *int,
) rangedIntField {
	return rangedIntField{intField{name, value}, lowest, highest}
}
