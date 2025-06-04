package game

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
	"gorm.io/gorm"
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
	Mu          sync.Mutex
	ID          int
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
	r.Mu.Lock()

	p.Room = r // Set the player's room
	r.Players[p.Name] = p
	msg := fmt.Sprintf(r.JoinMsg, p.Name)

	// !NOTE: Unlock before Broadcast to avoid deadlock
	r.Mu.Unlock()

	r.Broadcast(msg, p.Name)
}

func (r *Room) RemovePlayer(name string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.Players, name)

	msg := fmt.Sprintf(r.LeaveMsg, name)
	for otherName, player := range r.Players {
		if otherName != name {
			player.Conn.Write([]byte(utils.Colorize("\n"+msg+"\n", utils.Gray)))

			// Reprint prompt for interactivity
			player.Conn.Write([]byte(PlayerPrompt(player.Name, r.Name)))
		}
	}
}

func (r *Room) Broadcast(msg string, sender string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	for name, player := range r.Players {
		if name != sender {
			player.Conn.Write([]byte(utils.Colorize("\n"+msg+"\n", utils.Gray)))

			// Reprint prompt for interactivity
			player.Conn.Write([]byte(PlayerPrompt(player.Name, r.Name)))
		}
	}
}

func LoadAllRooms(gdb *gorm.DB) map[string]*Room {
	var roomModels []db.RoomModel
	if err := gdb.Find(&roomModels).Error; err != nil {
		log.Printf("Failed to load rooms: %v", err)
		return nil
	}

	// Create all Room instances with name keys
	roomMap := make(map[string]*Room)
	for _, model := range roomModels {
		room := &Room{
			ID:       int(model.ID),
			Name:     model.RoomName,
			Desc:     model.Description,
			Exits:    make(map[string]*Room),
			JoinMsg:  model.JoinMsg,
			LeaveMsg: model.LeaveMsg,
			Players:  make(map[string]*Player),
		}
		roomMap[room.Name] = room
	}

	// Wire-up Exits
	for _, model := range roomModels {
		cRoom := roomMap[model.RoomName]
		for _, exit := range model.Exits {
			parts := strings.SplitN(exit, ":", 2)
			if len(parts) == 2 {
				dir := strings.TrimSpace(parts[0])
				dName := strings.TrimSpace(parts[1])
				if dRoom, ok := roomMap[dName]; ok {
					cRoom.Exits[dir] = dRoom
				}
			}
		}

	}

	log.Printf("Loaded %d rooms", len(roomMap))

	return roomMap
}
