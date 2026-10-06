package main

import "sync"


type Hub struct {
	mu sync.Mutex  // prevents race conditions when accessing the clients map
	clients map[chan string]struct{}
}


// client subscription management
func (hub *Hub) Subscribe() chan string {
	channel := make(chan string, 32)  // 32 buffered channel
	
	hub.mu.Lock()
	hub.clients[channel] = struct{}{}  // dummy value
	hub.mu.Unlock()
	
	return channel
}


func (hub *Hub) Unsubscribe(channel chan string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	_, exists := hub.clients[channel]

	if !exists {
		return
	}

	delete(hub.clients, channel)
	close(channel)
}


func (hub *Hub) Broadcast(message string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for channel := range hub.clients {
		select {
			case channel <- message:
				// message sent successfully
			default:
				// disconnect a client which queue is full (not reading)
				delete(hub.clients, channel)
				close(channel)
		}
	}
}
