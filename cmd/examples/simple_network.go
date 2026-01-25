package main

import (
	"fmt"

	"github.com/lorenzhoerb/vroute/internal/algorithm"
	"github.com/lorenzhoerb/vroute/internal/router"
)

func main() {
	dijkstra := algorithm.NewDijkstra()

	a := router.NewRouter("a", dijkstra)
	b := router.NewRouter("b", dijkstra)
	c := router.NewRouter("c", dijkstra)
	d := router.NewRouter("d", dijkstra)

	router.Link(a, b, 1)
	router.Link(c, b, 2)
	router.Link(c, a, 5)
	router.Link(d, c, 3)

	// simulate flooding
	routers := []*router.Router{a, b, c, d}
	for _, r := range routers {
		for _, other := range routers {
			if r != other {
				other.ReceiveLSA(r.ID(), r.Neighbors())
			}
		}
	}

	// compute tables
	for _, r := range routers {
		if err := r.RecomputeRoutingTable(); err != nil {
			panic(err)
		}
	}

	for _, r := range routers {
		fmt.Printf("Routing table for %s:\n", r.ID())
		for dest, route := range r.RoutingTable() {
			fmt.Printf("  To %s via %s cost %.0f\n", dest, route.NextHop, route.Cost)
		}
		fmt.Println()
	}

}
