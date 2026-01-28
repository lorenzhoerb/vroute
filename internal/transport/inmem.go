package transport

import (
	"github.com/lorenzhoerb/vroute/internal/message"
	"github.com/lorenzhoerb/vroute/internal/topology"
)

type InMemTransport struct {
	id    topology.NodeID
	net   *InMemNetwork
	inbox chan message.Message
}

func NewInMemTransport(id topology.NodeID, net *InMemNetwork) *InMemTransport {
	return &InMemTransport{
		id:    id,
		net:   net,
		inbox: net.Register(id),
	}
}

func (t *InMemTransport) Send(to string, msg message.Message) error {
	t.net.Send(topology.NodeID(to), msg)
	return nil
}

func (t *InMemTransport) Receive() (<-chan message.Message, error) {
	return t.inbox, nil
}

func (t *InMemTransport) Close() {
	// nothing
}
