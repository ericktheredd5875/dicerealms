package game

import (
	"fmt"
	"log"
	"sync"
)

type Room struct {
	Name    string
	Players map[string]*Player
	mu      sync.Mutex
}

func NewRoom(name string) *Room {

	log.Printf("Creating new room: %s", name)
	return &Room{
		Name:    name,
		Players: make(map[string]*Player),
	}
}

func (r *Room) AddPlayer(p *Player) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in AddPlayer: %v", r)
		}
	}()

	r.mu.Lock()
	r.mu.Unlock() // !NOTE: Unlock before Broadcast to avoid deadlock

	log.Printf("Registering player: %s in room %s", p.Name, r.Name)
	r.Players[p.Name] = p
	log.Printf("Player registered in room: %s", r.Name)
	// r.Broadcast(fmt.Sprintf("%s has entered the room", p.Name), p.Name)
	message := fmt.Sprintf("%s has entered the room", p.Name)
	log.Printf("Message to broadcast: %s", message)
	r.Broadcast(message, p.Name)
	log.Printf("Broadcast sent to room: %s", r.Name)
}

func (r *Room) RemovePlayer(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Players, name)
}

func (r *Room) Broadcast(message string, sender string) {
	log.Printf("Starting broadcast in room: %s", r.Name)
	r.mu.Lock()
	log.Printf("Lock acquired for room: %s", r.Name)
	defer r.mu.Unlock()
	log.Printf("Lock released for room: %s", r.Name)

	log.Printf("Broadcasting in %s from %s to %d player(s)", r.Name, sender, len(r.Players))
	for name, player := range r.Players {
		if name != sender {
			player.Conn.Write([]byte(message + "\n"))
		}
	}
}
