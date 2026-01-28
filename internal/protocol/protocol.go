package protocol

import (
	"fmt"

	"github.com/lorenzhoerb/vroute/internal/message"
	"github.com/lorenzhoerb/vroute/internal/router"
	"github.com/lorenzhoerb/vroute/internal/topology"
)

// Protocol manages the routing logic and message exchange for a single router.
type Protocol struct {
	router    *router.Router
	transport Transport
}

// NewProtocol creates a new protocol instance for a router with a transport.
func NewProtocol(r *router.Router, t Transport) *Protocol {
	return &Protocol{
		router:    r,
		transport: t,
	}
}

// Run starts the protocol loop to process incoming messages.
func (p *Protocol) Run() error {
	ch, err := p.transport.Receive()
	if err != nil {
		return err
	}

	for msg := range ch {
		switch msg.Type {
		case message.LSA_Type:
			p.handleLSA(msg)
		case message.LSA_REQUEST_Type:
			p.handleLSARequest(msg)
		case message.DATA_Type:
			p.handleData(msg)
		default:
			fmt.Printf("[%s] unknown message type from %s\n", p.router.ID(), msg.Source)
		}
	}
	return nil
}

// SendData sends an arbitrary message to a destination router.
func (p *Protocol) SendData(dest string, data string) {
	msg := message.Message{
		Type:    message.DATA_Type,
		Source:  p.router.ID(),
		Dest:    topology.NodeID(dest),
		Payload: data,
	}
	p.routeMessage(msg)
}

// SendLSA floods the router's current LSA to all neighbors.
func (p *Protocol) SendLSA(lsa message.LSA) {
	msg := message.Message{
		Type:    message.LSA_Type,
		Source:  p.router.ID(),
		Payload: lsa,
	}

	p.flood(msg)
}

// RequestLSA asks all neighbors to send their current LSAs.
func (p *Protocol) RequestLSA() {
	msg := message.Message{
		Type:   message.LSA_REQUEST_Type,
		Source: p.router.ID(),
	}
	for neighbor := range p.router.Neighbors() {
		p.transport.Send(string(neighbor), msg)
	}
}

// handleData delivers or forwards messages based on the routing table.
func (p *Protocol) handleData(msg message.Message) {
	if p.router.ID() == msg.Dest {
		fmt.Printf("[%s] Received message from %s: %s\n", p.router.ID(), msg.Source, msg.Payload)
		return
	}

	p.routeMessage(msg)
}

// routeMessage forwards a message based on the router's routing table.
func (p *Protocol) routeMessage(msg message.Message) {
	// forward if dest is known
	routingTable := p.router.RoutingTable()
	route, ok := routingTable[topology.NodeID(msg.Dest)]
	if !ok {
		fmt.Printf("[%s] dropping msg from %s. Unknown destination: %s\n", p.router.ID(), msg.Source, msg.Dest)
		return
	}

	fmt.Printf("[%s] forwarding message to %s over %s\n", p.router.ID(), msg.Dest, route.NextHop)
	p.transport.Send(string(route.NextHop), msg)
}

// handleLSARequest responds to a neighbor's LSA request with our current LSA.
func (p *Protocol) handleLSARequest(msg message.Message) {
	lsa := p.router.GenerateLSA()
	lsaEchoMsg := message.Message{
		Type:    message.LSA_Type,
		Source:  p.router.ID(),
		Payload: lsa,
	}

	p.transport.Send(string(msg.Source), lsaEchoMsg)

}

// handleLSA processes an incoming LSA and floods it to neighbors if new.
func (p *Protocol) handleLSA(msg message.Message) {
	lsa, ok := msg.Payload.(message.LSA)
	if !ok {
		return
	}

	accept := p.router.ReceiveLSA(lsa)
	if !accept {
		// not accepted
		return
	}

	p.flood(msg)
}

// flood sends a message to all neighbors except the source.
func (p *Protocol) flood(msg message.Message) {
	for neighbor := range p.router.Neighbors() {
		if neighbor == msg.Source {
			continue
		}
		p.transport.Send(string(neighbor), msg)
	}
}
