package topology

type NodeID string

// Represents a node in a graph.
// It has a unique id and outgoing edges
type Node struct {
	ID    NodeID
	edges []*Edge
}

// Create a new node.
func NewNode(id NodeID) *Node {
	return &Node{ID: NodeID(id)}
}

// Adds an edge.
func (n *Node) AddEdge(to *Node, weight float64) {
	n.edges = append(n.edges, &Edge{
		To:     to,
		Weight: weight,
	})
}

// Edges gets all edges of the node.
func (n *Node) Edges() []*Edge {
	return n.edges
}

// Represents an Edge with a weight.
type Edge struct {
	Weight float64
	To     *Node
}
