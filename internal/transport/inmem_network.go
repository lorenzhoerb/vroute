package transport

import (
	"sync"

	"github.com/lorenzhoerb/vroute/internal/message"
	"github.com/lorenzhoerb/vroute/internal/topology"
)

type InMemNetwork struct {
	mu    sync.RWMutex
	nodes map[topology.NodeID]chan message.Message
}

func NewInMemNetwork() *InMemNetwork {
	return &InMemNetwork{
		nodes: make(map[topology.NodeID]chan message.Message),
	}
}

func (n *InMemNetwork) Register(id topology.NodeID) chan message.Message {
	n.mu.Lock()
	defer n.mu.Unlock()

	ch := make(chan message.Message, 16)
	n.nodes[id] = ch
	return ch
}

func (n *InMemNetwork) Send(to topology.NodeID, msg message.Message) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if ch, ok := n.nodes[to]; ok {
		ch <- msg
	}
}
