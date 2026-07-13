package misc

type pairSet map[[2]int]bool

//nolint:revive // This is just a type alias for a map, I don't want this type to be referenced
func NewPairSet() *pairSet {
	set := make(pairSet)
	return &set
}

func (this *pairSet) Add(a, b int) {
	if a > b {
		a, b = b, a
	}
	(*this)[[2]int{a, b}] = true
}
