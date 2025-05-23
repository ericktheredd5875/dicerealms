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
	r.mu.Lock()
	r.mu.Unlock() // !NOTE: Unlock before Broadcast to avoid deadlock

	p.Room = r // Set the player's room
	r.Players[p.Name] = p
	r.Broadcast(fmt.Sprintf("%s has entered the room", p.Name), p.Name)
}

func (r *Room) RemovePlayer(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Players, name)
}

func (r *Room) Broadcast(message string, sender string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, player := range r.Players {
		if name != sender {
			player.Conn.Write([]byte(message + "\n"))
		}
	}
}
