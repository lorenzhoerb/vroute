package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph_AddNode(t *testing.T) {
	tests := []struct {
		name        string
		setupGraph  func() *Graph
		inputNodeId NodeID
		wantNode    *Node
	}{
		{
			name: "creates new node",
			setupGraph: func() *Graph {
				return NewGraph()
			},
			inputNodeId: "n1",
			wantNode:    nil, // nil means a new node is expected
		},
		{
			name: "returns existing node",
			setupGraph: func() *Graph {
				return &Graph{
					nodes: map[NodeID]*Node{
						"n1": {ID: "n1"},
					},
				}

			},
			inputNodeId: "n1",
			wantNode:    nil, // pointer comparison will be done after add
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.setupGraph()

			prevNode := g.nodes[tt.inputNodeId] // direct access

			gotNode := g.AddNode(tt.inputNodeId)

			assert.NotNil(t, gotNode, "node should not be nil")
			assert.Equal(t, tt.inputNodeId, gotNode.ID, "node ID mismatch")

			if prevNode != nil {
				assert.Same(t, prevNode, gotNode, "should return the same existing node")
			}
		})
	}
}
