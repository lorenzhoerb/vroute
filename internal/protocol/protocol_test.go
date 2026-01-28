package protocol

import (
	"fmt"
	"testing"
	"time"

	"github.com/lorenzhoerb/vroute/internal/algorithm"
	"github.com/lorenzhoerb/vroute/internal/router"
	"github.com/lorenzhoerb/vroute/internal/transport"
)

// Topology
// A-B-C-D
//
//	\  /
//
// Send message from A to D
// Shortest Path is A-B-D
func TestProtocol(t *testing.T) {
	net := transport.NewInMemNetwork()

	algo := algorithm.NewDijkstra()

	ra := router.NewRouter("A", algo)
	rb := router.NewRouter("B", algo)
	rc := router.NewRouter("C", algo)
	rd := router.NewRouter("D", algo)

	ta := transport.NewInMemTransport("A", net)
	tb := transport.NewInMemTransport("B", net)
	tc := transport.NewInMemTransport("C", net)
	td := transport.NewInMemTransport("D", net)

	pa := Protocol{ra, ta}
	pb := Protocol{rb, tb}
	pc := Protocol{rc, tc}
	pd := Protocol{rd, td}

	go pa.Run()
	go pb.Run()
	go pc.Run()
	go pd.Run()

	// link A <-> B (ID only)
	ra.UpdateNeighbor("B", 1)
	lsaA := ra.GenerateLSA()
	pa.SendLSA(lsaA)
	pa.RequestLSA()

	rc.UpdateNeighbor("B", 2)
	lsaC := rc.GenerateLSA()
	pc.SendLSA(lsaC)
	pc.RequestLSA()

	rc.UpdateNeighbor("D", 5)
	lsaC = rc.GenerateLSA()
	pc.SendLSA(lsaC)
	pc.RequestLSA()

	rd.UpdateNeighbor("B", 1)
	lsaD := rd.GenerateLSA()
	pd.SendLSA(lsaD)
	pd.RequestLSA()

	time.Sleep(200 * time.Millisecond)
	pa.router.Recompute()
	pb.router.Recompute()
	pc.router.Recompute()
	pd.router.Recompute()

	fmt.Println("Router A - Routing Table")
	fmt.Print(ra.RoutingTable())
	fmt.Println()

	fmt.Println("Router B - Routing Table")
	fmt.Print(rb.RoutingTable())
	fmt.Println()

	fmt.Println("Router C - Routing Table")
	fmt.Print(rc.RoutingTable())
	fmt.Println()

	fmt.Println("Router D - Routing Table")
	fmt.Print(rd.RoutingTable())
	fmt.Println()

	pa.SendData("D", "hello world")

}
