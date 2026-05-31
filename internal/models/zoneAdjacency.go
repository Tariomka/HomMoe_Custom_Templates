package models

type ZoneAdjacency map[string]map[string]bool

func (this *ZoneAdjacency) Link(inputLabel, outputLabel string) {
	if (*this)[inputLabel] == nil {
		(*this)[inputLabel] = map[string]bool{}
	}
	if (*this)[outputLabel] == nil {
		(*this)[outputLabel] = map[string]bool{}
	}
	(*this)[inputLabel][outputLabel] = true
	(*this)[outputLabel][inputLabel] = true
}

func (this *ZoneAdjacency) GetDistancesFrom(startLabel string) map[string]int {
	distances := map[string]int{startLabel: 0}
	queue := []string{startLabel}
	for len(queue) > 0 {
		currentLabel := queue[0]
		queue = queue[1:]
		for nextLabel := range (*this)[currentLabel] {
			if _, ok := distances[nextLabel]; !ok {
				distances[nextLabel] = distances[currentLabel] + 1
				queue = append(queue, nextLabel)
			}
		}
	}
	return distances
}
