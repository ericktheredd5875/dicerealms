package game

import (
	"fmt"
)

func (p *Player) Narrate(msg string) {
	if msg == "" {
		return
	}

	message := fmt.Sprintf("<narritive> %s", msg)
	p.Room.Broadcast(message, "")

	if p.Room.ActiveScene != nil {
		p.Room.ActiveScene.LogEntry(message)
	}
}
