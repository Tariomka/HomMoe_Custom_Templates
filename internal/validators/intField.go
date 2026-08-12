package validators

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

// intField names a .gen.json integer field and provides access to it, so a
// single check can both read the value and apply its fix.
type intField struct {
	name  string
	value func(state *dtos.EditorStateDto) *int
}
