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

type ZoneIndexAdjacency map[int]map[int]bool

func NewZoneIndexAdjacency(size int) *ZoneIndexAdjacency {
	adjacency := make(ZoneIndexAdjacency, size)
	for i := range adjacency {
		adjacency[i] = make(map[int]bool)
	}
	return &adjacency
}

func (this *ZoneIndexAdjacency) Link(inputIdx, outputIdx int) {
	(*this)[inputIdx][outputIdx] = true
	(*this)[outputIdx][inputIdx] = true
}

func (this *ZoneIndexAdjacency) FindIndexes(nodeCount int) [][]int {
	visited := make([]bool, nodeCount)
	var components [][]int
	for start := range nodeCount {
		if visited[start] {
			continue
		}
		var comp []int
		queue := []int{start}
		visited[start] = true
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			comp = append(comp, cur)
			for nb := range (*this)[cur] {
				if !visited[nb] {
					visited[nb] = true
					queue = append(queue, nb)
				}
			}
		}
		components = append(components, comp)
	}
	return components
}
