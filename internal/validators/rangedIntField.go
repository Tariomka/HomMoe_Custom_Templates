package validators

type rangedIntField struct {
	field   intField
	lowest  int
	highest int
}
