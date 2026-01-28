package protocol

import "github.com/lorenzhoerb/vroute/internal/message"

type Transport interface {
	Send(to string, message message.Message) error
	Receive() (<-chan message.Message, error)
	Close()
}
