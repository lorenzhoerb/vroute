package router

import (
	"fmt"
	"strings"

	"github.com/lorenzhoerb/vroute/internal/topology"
)

type Route struct {
	Destination topology.NodeID
	NextHop     topology.NodeID
	Cost        float64
}

type RoutingTale map[topology.NodeID]Route

func (r *Router) RoutingTable() RoutingTale {
	return r.routingTale
}

func (r RoutingTale) String() string {
	if len(r) == 0 {
		return "<empty>"
	}

	var sb strings.Builder
	for dest, route := range r {
		sb.WriteString(fmt.Sprintf("To %s via %s cost %.2f\n", dest, route.NextHop, route.Cost))
	}
	return sb.String()
}
