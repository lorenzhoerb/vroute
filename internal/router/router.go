package router

import (
	"github.com/lorenzhoerb/vroute/internal/algorithm"
	"github.com/lorenzhoerb/vroute/internal/topology"
)

type Router struct {
	id topology.NodeID

	neighbors map[topology.NodeID]float64

	lsdb map[topology.NodeID]map[topology.NodeID]float64

	routingTale RoutingTale

	// derived fom lsbd
	graph *topology.Graph

	algo algorithm.PathFinder
}

func NewRouter(id topology.NodeID, algo algorithm.PathFinder) *Router {
	return &Router{
		id:        id,
		neighbors: make(map[topology.NodeID]float64),
		lsdb: map[topology.NodeID]map[topology.NodeID]float64{
			id: {},
		},
		algo: algo,
	}
}

func (r *Router) ID() topology.NodeID {
	return r.id
}

func (r *Router) Neighbors() map[topology.NodeID]float64 {
	return r.neighbors
}

func (r *Router) UpdateNeighbor(id topology.NodeID, cost float64) {
	r.neighbors[id] = cost
	r.lsdb[r.id] = r.neighbors
}

func (r *Router) DeleteNeighbor(id topology.NodeID) {
	delete(r.neighbors, id)
	r.lsdb[r.id] = r.neighbors
}

func (r *Router) ReceiveLSA(from topology.NodeID, neighbors map[topology.NodeID]float64) {
	r.lsdb[from] = neighbors
}

func (r *Router) rebuildGraph() {
	g := topology.NewGraph()

	for from, neighs := range r.lsdb {
		for to, cost := range neighs {
			g.AddEdge(from, to, cost)
		}
	}

	r.graph = g
}
