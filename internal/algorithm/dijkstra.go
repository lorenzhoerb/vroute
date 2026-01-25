package algorithm

import (
	"fmt"
	"math"

	"github.com/lorenzhoerb/vroute/internal/topology"
)

type dijkstra struct{}

func NewDijkstra() *dijkstra {
	return &dijkstra{}
}

func (d *dijkstra) ComputeShortestPaths(graph *topology.Graph, source topology.NodeID) (ShortestPaths, error) {
	shortestPaths := make(ShortestPaths)

	sourceNode, ok := graph.GetNode(source)
	if !ok {
		return nil, fmt.Errorf("source node %s does not exist", source)
	}

	// initialize all nodes with infinite cost and nil previous
	for _, node := range graph.Nodes() {
		shortestPaths[node.ID] = PathInfo{
			Cost:     math.Inf(1),
			PrevNode: nil,
		}
	}

	// init source node in shortest path
	shortestPaths[sourceNode.ID] = PathInfo{Cost: 0, PrevNode: nil}

	visited := make(map[topology.NodeID]bool)

	for {
		// select the unvisited node with the smallest cost
		var currentNode *topology.Node
		minCost := math.Inf(1)

		for nodeID, info := range shortestPaths {
			if !visited[nodeID] && info.Cost < minCost {
				minCost = info.Cost
				currentNode, _ = graph.GetNode(nodeID)
			}
		}

		if currentNode == nil {
			break // no more reachable nodes
		}

		visited[currentNode.ID] = true

		for _, edge := range currentNode.Edges() {
			if visited[edge.To.ID] {
				continue
			}
			newCost := shortestPaths[currentNode.ID].Cost + edge.Weight
			if newCost < shortestPaths[edge.To.ID].Cost {
				shortestPaths[edge.To.ID] = PathInfo{
					Cost:     newCost,
					PrevNode: currentNode,
				}
			}
		}
	}

	return shortestPaths, nil
}
