package message

import "github.com/lorenzhoerb/vroute/internal/topology"

type MessageType uint8

const (
	LSA_Type MessageType = iota
	LSA_REQUEST_Type
	DATA_Type
)

type Message struct {
	Type    MessageType
	Source  topology.NodeID
	Dest    topology.NodeID
	Payload any
}

type LSA struct {
	Origin topology.NodeID

	// monotonically increasing per origin router
	Sequence uint64

	// adjacency list
	Links map[topology.NodeID]float64
}
