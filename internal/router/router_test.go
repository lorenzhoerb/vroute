package router

import (
	"testing"

	"github.com/lorenzhoerb/vroute/internal/algorithm"
	"github.com/lorenzhoerb/vroute/internal/message"
	"github.com/lorenzhoerb/vroute/internal/topology"
	"github.com/stretchr/testify/assert"
)

func TestRoutingExample(t *testing.T) {
	dijkstra := algorithm.NewDijkstra()

	a := NewRouter("A", dijkstra)
	b := NewRouter("B", dijkstra)
	c := NewRouter("C", dijkstra)
	d := NewRouter("D", dijkstra)

	// establish links
	Link(a, b, 1)
	Link(b, c, 2)
	Link(a, c, 5)
	Link(c, d, 3)

	routers := []*Router{a, b, c, d}

	// simulate flooding
	for _, r := range routers {
		for _, other := range routers {
			if r != other {
				other.ReceiveLSA(message.LSA{
					Origin:   r.id,
					Sequence: 0,
					Links:    r.Neighbors(),
				})
			}
		}
	}

	// compute routing tables
	for _, r := range routers {
		if err := r.Recompute(); err != nil {
			t.Fatal(err)
		}
	}

	// assert routing tables
	assert.Equal(t, topology.NodeID("B"), a.RoutingTable()["B"].NextHop)
	assert.Equal(t, topology.NodeID("B"), a.RoutingTable()["C"].NextHop)
	assert.Equal(t, topology.NodeID("B"), a.RoutingTable()["D"].NextHop)
}
