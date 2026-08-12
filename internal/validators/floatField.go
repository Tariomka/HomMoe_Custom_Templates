package validators

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

// floatField names a .gen.json float field and provides access to it, so a
// single check can both read the value and apply its fix.
type floatField struct {
	name  string
	value func(state *dtos.EditorStateDto) *float64
}
