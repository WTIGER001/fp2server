package common

import (
	"context"
	"sync"
)

type MessageBus interface {
	Send(ctx context.Context, message *Fp2Message) error
	RegisterHandler(h MessageHandler)
	UnregisterHandler(h MessageHandler)
}

// Global Reference to the Comms service
var Comms MessageBus

type MessageHandler interface {
	Handle(message *Fp2Message)
}

type SingleResponseHandler struct {
	MessageBus MessageBus
	Message    *Fp2Message
	Response   *Fp2Message
	WG         *sync.WaitGroup
	OnRecieve  func(message *Fp2Message)
	OnTimeout  func(message *Fp2Message)
}

func (h SingleResponseHandler) Handle(message *Fp2Message) {
	if message.RespondingToID == h.Message.MessageID {
		h.Response = message
		if h.OnRecieve != nil {
			h.OnRecieve(message)
		}
		h.MessageBus.UnregisterHandler(h)
		h.WG.Done()
	}
}

func (h SingleResponseHandler) Send(ctx context.Context) {
	h.WG.Add(1)
	h.MessageBus.RegisterHandler(h)
	h.MessageBus.Send(ctx, h.Message)
}

func prepMessage(m *Fp2Message) {
	// Verify there is an id
	if m.MessageID == "" {
		m.MessageID = GenerateID()
	}
	// Verify there is a sender
	m.Sender = "_SERVER_"
}

//Sends out a message
func SendAndWait(ctx context.Context, m *Fp2Message) *Fp2Message {
	prepMessage(m)

	// Create the Handlers
	h := &SingleResponseHandler{
		MessageBus: Comms,
		Message:    m,
		WG:         new(sync.WaitGroup),
	}
	h.Send(ctx)

	// Wait
	h.WG.Wait()

	return h.Response
}

func SendAndCallback(ctx context.Context, m *Fp2Message, fn func(message *Fp2Message)) {
	prepMessage(m)

	// Create the Handlers
	h := &SingleResponseHandler{
		MessageBus: Comms,
		Message:    m,
		WG:         new(sync.WaitGroup),
		OnRecieve:  fn,
	}
	h.Send(ctx)
}

// func (cc *CommChain) Handle(message *Fp2Message) {

// }

// // Send a bunch of messages and wait for them all to be resolved
// // Basically fork join
// func (cc *CommChain) SendAllAndWait() {

// }

// // Resolve a timeout
// func (cc *CommChain) ResolveTimeout() {

// }
