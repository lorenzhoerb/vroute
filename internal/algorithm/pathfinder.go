package algorithm

import "github.com/lorenzhoerb/vroute/internal/topology"

// PathInfo stores the cost from the source and the previous node
type PathInfo struct {
	Cost     float64
	PrevNode *topology.Node
}

// ShortestPaths maps node IDs to PathInfo
type ShortestPaths map[topology.NodeID]PathInfo

// PathFinder computes shortest paths from a source node
type PathFinder interface {
	ComputeShortestPaths(graph *topology.Graph, source topology.NodeID) (ShortestPaths, error)
}
