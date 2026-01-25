package topology

// Graph represents a graph.
type Graph struct {
	nodes map[NodeID]*Node
}

// NewGraph create a graph instance.
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[NodeID]*Node),
	}
}

// AddNode returns an existing node or creates a new one.
func (g *Graph) AddNode(id NodeID) *Node {
	if n, ok := g.nodes[id]; ok {
		return n
	}

	n := NewNode(id)
	g.nodes[id] = n
	return n
}

// GetNodes returns the node if it exists.
// If exists returns (node, true),
// if not (nil, false).
func (g *Graph) GetNode(id NodeID) (*Node, bool) {
	n, ok := g.nodes[id]
	return n, ok
}

// AddEdge adds an edge from node `from` to node `to` with the given weight.
// It ensures that both nodes exist.
func (g *Graph) AddEdge(from, to NodeID, weight float64) {
	fromNode := g.AddNode(from)
	toNode := g.AddNode(to)

	fromNode.AddEdge(toNode, weight)
}

// AddEdge adds an bidirectional edge between nodes `a` and `b` with the given weight.
// It ensures that both nodes exist.
func (g *Graph) AddBidirectionalEdge(a, b NodeID, weight float64) {
	g.AddEdge(a, b, weight)
	g.AddEdge(b, a, weight)
}

func (g *Graph) Nodes() []*Node {
	nodes := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}
