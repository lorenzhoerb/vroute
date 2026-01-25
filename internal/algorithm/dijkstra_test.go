package algorithm

import (
	"math"
	"testing"

	"github.com/lorenzhoerb/vroute/internal/topology"
	"github.com/stretchr/testify/assert"
)

func TestDijkstra_ComputeShortestPaths(t *testing.T) {
	// Helper to compare floats with tolerance
	equalFloat := func(a, b float64) bool {
		const eps = 1e-9
		return math.Abs(a-b) < eps
	}

	t.Run("source does not exist", func(t *testing.T) {
		graph := topology.NewGraph()
		_, err := NewDijkstra().ComputeShortestPaths(graph, "X")
		assert.Error(t, err)
	})

	t.Run("simple graph", func(t *testing.T) {
		graph := topology.NewGraph()

		// Build graph:
		// A --1--> B --2--> C
		// A --4--> C
		graph.AddEdge("A", "B", 1)
		graph.AddEdge("B", "C", 2)
		graph.AddEdge("A", "C", 4)

		sp, err := NewDijkstra().ComputeShortestPaths(graph, "A")
		assert.NoError(t, err)

		// Expected costs
		expected := map[topology.NodeID]float64{
			"A": 0,
			"B": 1,
			"C": 3, // via B, cheaper than direct 4
		}

		for nodeID, cost := range expected {
			info, ok := sp[nodeID]
			assert.True(t, ok, "node %s should exist in shortest paths", nodeID)
			assert.True(t, equalFloat(cost, info.Cost), "node %s expected cost %.2f got %.2f", nodeID, cost, info.Cost)
		}

		// Check PrevNode chain for C
		infoC := sp["C"]
		assert.Equal(t, topology.NodeID("B"), infoC.PrevNode.ID)
		infoB := sp["B"]
		assert.Equal(t, topology.NodeID("A"), infoB.PrevNode.ID)
	})

	t.Run("graph with unreachable node", func(t *testing.T) {
		graph := topology.NewGraph()
		graph.AddNode("A")
		graph.AddNode("B") // disconnected

		sp, err := NewDijkstra().ComputeShortestPaths(graph, "A")
		assert.NoError(t, err)

		assert.Equal(t, 0.0, sp["A"].Cost)
		assert.Equal(t, math.Inf(1), sp["B"].Cost) // unreachable node
		assert.Nil(t, sp["B"].PrevNode)
	})
}
