package graph

type Adjacency[Node comparable] map[Node]map[Node]bool

func NewAdjacency[Node comparable](nodes []Node) Adjacency[Node] {
	adjacency := make(Adjacency[Node], len(nodes))
	for _, node := range nodes {
		adjacency[node] = make(map[Node]bool)
	}
	return adjacency
}

func Link[Node comparable](adjacency Adjacency[Node], inputNode, outputNode Node) {
	if adjacency[inputNode] == nil {
		adjacency[inputNode] = map[Node]bool{}
	}
	if adjacency[outputNode] == nil {
		adjacency[outputNode] = map[Node]bool{}
	}
	adjacency[inputNode][outputNode] = true
	adjacency[outputNode][inputNode] = true
}

func DistancesFrom[Node comparable](adjacency Adjacency[Node], startNode Node) map[Node]int {
	distances := map[Node]int{startNode: 0}
	queue := []Node{startNode}
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		for nextNode := range adjacency[currentNode] {
			if _, visited := distances[nextNode]; !visited {
				distances[nextNode] = distances[currentNode] + 1
				queue = append(queue, nextNode)
			}
		}
	}
	return distances
}

func ConnectedComponents[Node comparable](adjacency Adjacency[Node], nodes []Node) [][]Node {
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
			for nextNode := range adjacency[currentNode] {
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
