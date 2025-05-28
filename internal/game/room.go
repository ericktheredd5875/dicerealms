package game

import (
	"fmt"
	"sync"
	"time"
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
	r.mu.Unlock() // !NOTE: Unlock before Broadcast to avoid deadlock

	p.Room = r // Set the player's room
	r.Players[p.Name] = p
	msg := fmt.Sprintf(r.JoinMsg, p.Name)
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
			player.Conn.Write([]byte(message + "\n"))
			// player.Conn.Write([]byte("+>> "))
		}
	}
}

func (r *Room) StartScene(title string, mood, startedBy string) {
	r.mu.Lock()
	r.mu.Unlock()

	r.ActiveScene = &Scene{
		Title:     title,
		Mood:      mood,
		StartedBy: startedBy,
		StartedAt: time.Now(),
		Log:       []string{},
	}

	r.Broadcast(fmt.Sprintf("<new> Scene started: %s [%s]", title, mood), "")
}

func (r *Room) EndScene(endedBy string) string {
	r.mu.Lock()
	r.mu.Unlock()

	if r.ActiveScene == nil {
		return "No active scene to end."
	}

	r.ActiveScene.EndedBy = endedBy
	r.ActiveScene.EndedAt = time.Now()
	summary := r.ActiveScene.Summary()

	r.ActiveScene = nil
	r.Broadcast(fmt.Sprintf("<end> Scene ended: %s", summary), "")

	return summary

}
