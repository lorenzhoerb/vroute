package router

import "github.com/lorenzhoerb/vroute/internal/topology"

type Route struct {
	Destination topology.NodeID
	NextHop     topology.NodeID
	Cost        float64
}

type RoutingTale map[topology.NodeID]Route

func (r *Router) RecomputeRoutingTable() error {
	r.rebuildGraph()

	paths, err := r.algo.ComputeShortestPaths(r.graph, r.id)
	if err != nil {
		return err
	}

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

	return nil
}

func (r *Router) RoutingTable() RoutingTale {
	return r.routingTale
}
