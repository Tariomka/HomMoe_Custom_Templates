package misc

type pairSet map[[2]int]bool

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
