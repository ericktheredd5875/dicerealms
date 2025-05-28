package game

import (
	"fmt"
	"sync"
)

type Room struct {
	Name    string
	Desc    string
	Players map[string]*Player
	// **IE: North, South, East, West, Up, Down, In, Out, etc.
	Exits       map[string]*Room
	JoinMsg     string
	LeaveMsg    string
	ActiveScene *Scene
	mu          sync.Mutex
}

func NewRoom(name string, desc string) *Room {
	return &Room{
		Name:     name,
		Desc:     desc,
		Players:  make(map[string]*Player),
		Exits:    make(map[string]*Room),
		JoinMsg:  "%s has entered the room.",
		LeaveMsg: "%s has left the room.",
	}
}

func (r *Room) AddPlayer(p *Player) {
	r.mu.Lock()

	p.Room = r // Set the player's room
	r.Players[p.Name] = p
	msg := fmt.Sprintf(r.JoinMsg, p.Name)

	// !NOTE: Unlock before Broadcast to avoid deadlock
	r.mu.Unlock()

	r.Broadcast(msg, p.Name)
}

func (r *Room) RemovePlayer(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Players, name)

	msg := fmt.Sprintf(r.LeaveMsg, name)
	for otherName, player := range r.Players {
		if otherName != name {
			player.Conn.Write([]byte(msg + "\n"))
			// player.Conn.Write([]byte("+>> "))
		}
	}
}

func (r *Room) Broadcast(message string, sender string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, player := range r.Players {
		if name != sender {
			player.Conn.Write([]byte("\n" + message + "\n"))

			// Reprint prompt for interactivity
			prompt := fmt.Sprintf("\n%s@%s +>> ", player.Name, r.Name)
			player.Conn.Write([]byte(prompt))
		}
	}
}
