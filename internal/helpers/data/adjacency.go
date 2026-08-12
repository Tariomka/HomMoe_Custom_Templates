package data

type Adjacency[Node comparable] map[Node]map[Node]bool

func NewAdjacency[Node comparable](nodes []Node) Adjacency[Node] {
	adjacency := make(Adjacency[Node])
	for _, node := range nodes {
		adjacency[node] = make(map[Node]bool)
	}
	return adjacency
}

func (this Adjacency[Node]) Link(inputNode, outputNode Node) {
	if this[inputNode] == nil {
		this[inputNode] = map[Node]bool{}
	}
	if this[outputNode] == nil {
		this[outputNode] = map[Node]bool{}
	}
	this[inputNode][outputNode] = true
	this[outputNode][inputNode] = true
}

func (this Adjacency[Node]) DistancesFrom(startNode Node) map[Node]int {
	distances := map[Node]int{startNode: 0}
	queue := []Node{startNode}
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		for nextNode := range this[currentNode] {
			if _, visited := distances[nextNode]; !visited {
				distances[nextNode] = distances[currentNode] + 1
				queue = append(queue, nextNode)
			}
		}
	}
	return distances
}

func (this Adjacency[Node]) ConnectedComponents(nodes []Node) [][]Node {
	visited := make(map[Node]bool, len(nodes))
	components := make([][]Node, 0)
	for _, startNode := range nodes {
		if visited[startNode] {
			continue
		}
		component := make([]Node, 0)
		queue := []Node{startNode}
		visited[startNode] = true
		for len(queue) > 0 {
			currentNode := queue[0]
			queue = queue[1:]
			component = append(component, currentNode)
			for nextNode := range this[currentNode] {
				if !visited[nextNode] {
					visited[nextNode] = true
					queue = append(queue, nextNode)
				}
			}
		}
		components = append(components, component)
	}
	return components
}
