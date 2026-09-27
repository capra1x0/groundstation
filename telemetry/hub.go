package main

import "sync"

type Hub struct {
	mu sync.Mutex
	clients map[chan Reading]bool
}

func newHub() *Hub {
	return &Hub{clients: map[chan Reading]bool{}}
}

func (h *Hub) add() chan Reading {
	ch := make(chan Reading, 64)

	h.mu.Lock()
	h.clients[ch] = true
	h.mu.Unlock()

	return ch
}

func (h *Hub) remove(ch chan Reading) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
}

func (h *Hub) broadcase(reading Reading) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for ch := range h.clients {
		select {
		case ch <- reading:
		default:

		}
	}
}
