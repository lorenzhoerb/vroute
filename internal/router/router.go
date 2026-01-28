package router

import (
	"github.com/lorenzhoerb/vroute/internal/algorithm"
	"github.com/lorenzhoerb/vroute/internal/message"
	"github.com/lorenzhoerb/vroute/internal/topology"
)

type LSAEntry struct {
	Sequence uint64
	Links    map[topology.NodeID]float64
}

// Router represents a network router.
type Router struct {
	id topology.NodeID

	neighbors map[topology.NodeID]float64
	selfSeq   uint64

	lsdb      map[topology.NodeID]LSAEntry
	lsdbDirty bool // True if lsdb has unprocessed changes

	routingTale RoutingTale // current routing table

	graph *topology.Graph // derived topology from lsdb

	algo algorithm.PathFinder // Pathfinding algorithm
}

func NewRouter(id topology.NodeID, algo algorithm.PathFinder) *Router {
	return &Router{
		id:        id,
		neighbors: make(map[topology.NodeID]float64),
		lsdb: map[topology.NodeID]LSAEntry{
			id: {
				Sequence: 0,
				Links:    make(map[topology.NodeID]float64),
			},
		},
		algo: algo,
	}
}

// ID returns the id of the router.
func (r *Router) ID() topology.NodeID {
	return r.id
}

// Neighbors returns the neighbors of the router.
func (r *Router) Neighbors() map[topology.NodeID]float64 {
	return r.neighbors
}

// Update Neighbor, updates the cost of a given neighbor.
func (r *Router) UpdateNeighbor(id topology.NodeID, cost float64) {
	r.neighbors[id] = cost
	r.updateSelfLSA()
}

// Delete Neighbors deletes the given neighbor and update the LSA.
func (r *Router) DeleteNeighbor(id topology.NodeID) {
	delete(r.neighbors, id)
	r.updateSelfLSA()
}

func (r *Router) updateSelfLSA() {
	r.selfSeq++

	r.lsdb[r.id] = LSAEntry{
		Sequence: r.selfSeq,
		Links:    r.neighbors,
	}
	r.lsdbDirty = true
}

// ReceiveLSA processes an incoming LSA (Link-State Advertisement) from another router.
// - If the LSA is older or already seen (sequence number <= stored), it is ignored and returns false.
// - If the LSA is newer, it updates the LSDB with the new links and marks the LSDB as dirty.
// - Additionally, if this router is listed as a neighbor in the LSA, it updates its own neighbor table.
// Returns true if the LSA was accepted and applied.
func (r *Router) ReceiveLSA(lsa message.LSA) bool {
	entry, exists := r.lsdb[lsa.Origin]

	if exists && lsa.Sequence <= entry.Sequence {
		return false
	}

	r.lsdb[lsa.Origin] = LSAEntry{
		Sequence: lsa.Sequence,
		Links:    lsa.Links,
	}

	// If this router is included in the LSA links, update neighbor cost
	if cost, ok := lsa.Links[r.id]; ok {
		r.neighbors[lsa.Origin] = cost
	}

	r.lsdbDirty = true
	return true
}

func (r *Router) GenerateLSA() message.LSA {
	return message.LSA{
		Origin:   r.id,
		Sequence: r.selfSeq,
		Links:    r.neighbors,
	}
}

// Recomputes recomputes the the routing table, if lsdb was changed.
func (r *Router) Recompute() error {
	if !r.lsdbDirty {
		return nil
	}

	r.rebuildGraph()

	paths, err := r.algo.ComputeShortestPaths(r.graph, r.id)
	if err != nil {
		return err
	}

	r.buildRoutingTable(paths)

	r.lsdbDirty = false
	return nil

}

func (r *Router) buildRoutingTable(paths algorithm.ShortestPaths) {
	table := make(RoutingTale)

	for dest, info := range paths {
		if dest == r.id || info.PrevNode == nil {
			continue
		}

		nextHop := dest
		for paths[nextHop].PrevNode.ID != r.id {
			nextHop = paths[nextHop].PrevNode.ID
		}

		table[dest] = Route{
			Destination: dest,
			NextHop:     nextHop,
			Cost:        info.Cost,
		}

		r.routingTale = table
	}
}

func (r *Router) rebuildGraph() {
	g := topology.NewGraph()

	for from, lsaEntry := range r.lsdb {
		for to, cost := range lsaEntry.Links {
			g.AddBidirectionalEdge(from, to, cost)
		}
	}

	r.graph = g
}
