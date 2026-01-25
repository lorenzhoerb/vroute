package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	id := NodeID("node-id")
	n := NewNode(id)

	assert.NotNil(t, n, "node should not be nil")
	assert.Equal(t, id, n.ID, "node ID should match")
}

func TestNode_AddEdge(t *testing.T) {
	t.Run("add single edge", func(t *testing.T) {
		n1 := NewNode(NodeID("n1"))
		n2 := NewNode(NodeID("n2"))
		w := 1.0

		n1.AddEdge(n2, w)

		assert.Len(t, n1.edges, 1, "edge slice of n1 should have length of 1")
		assert.Len(t, n2.edges, 0, "edge slice of n2 should have length of 0")

		assert.Equal(t, n2, n1.edges[0].To, "only edge in n1.edges should point to n2")
		assert.Equal(t, w, n1.edges[0].Weight, "only edge in n1.edges should have weight of 1")

	})

	t.Run("add multiple edges", func(t *testing.T) {
		n1 := NewNode(NodeID("n1"))
		n2 := NewNode(NodeID("n2"))
		n3 := NewNode(NodeID("n3"))

		n1.AddEdge(n2, 1.0)
		n1.AddEdge(n3, 2.0)

		assert.Len(t, n1.edges, 2, "edge slice of n1 should have length of 3")
		assert.Len(t, n2.edges, 0, "edge slice of n2 should have length of 0")
		assert.Len(t, n3.edges, 0, "edge slice of n3 should have length of 0")

		assert.Equal(t, n2, n1.edges[0].To, "first edge should go to n2")
		assert.Equal(t, n3, n1.edges[1].To, "second edge should go to n3")
		assert.Equal(t, 1.0, n1.edges[0].Weight, "first edge should have wight of 1.0")
		assert.Equal(t, 2.0, n1.edges[1].Weight, "first edge should have weight of 2.0")

	})
}

func TestNode_Edges(t *testing.T) {
	n1 := NewNode(NodeID("n1"))
	n2 := NewNode(NodeID("n2"))
	n3 := NewNode(NodeID("n3"))

	n1.AddEdge(n2, 1.0)
	n1.AddEdge(n3, 2.0)

	gotEdges := n1.Edges()

	assert.Len(t, gotEdges, 2, "should return two edges")
	assert.Equal(t, n2, gotEdges[0].To, "first edge should go to n2")
	assert.Equal(t, n3, gotEdges[1].To, "second edge should go to n3")
}
