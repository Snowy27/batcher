package models

type node struct {
	Visited      bool
	Verified     bool
	Dependencies []string
}

//DependenciesGraph represents the requests dependencies
type DependenciesGraph map[string]*node

//NewDependenciesGraph creates a DependenciesGraph from array of requests
func NewDependenciesGraph(requests []Request) DependenciesGraph {
	graph := make(DependenciesGraph)
	for _, request := range requests {
		graph[request.Name] = &node{Dependencies: request.Dependencies}
	}
	return graph
}

//CheckForCircularDependencies traverses the graph and returns true if there are circular dependencies, otherwise returns false
func (graph DependenciesGraph) CheckForCircularDependencies() bool {
	for _, node := range graph {
		if graph.checkNode(node) {
			return true
		}
	}

	return false
}

//checkNode recursive function that traverses each node's dependencies and returns true if it encounters circular dependency
func (graph DependenciesGraph) checkNode(node *node) bool {
	if node.Verified {
		return false
	}

	if node.Visited {
		return true
	}

	node.Visited = true

	for _, dependency := range node.Dependencies {
		if graph.checkNode(graph[dependency]) {
			return true
		}
	}

	node.Verified = true
	return false
}
