package game

import (
	"fmt"
)

func (p *Player) Move(dir string) (string, error) {
	p.Room.mu.Lock()
	targetRoom, ok := p.Room.Exits[dir]
	p.Room.mu.Unlock()

	if !ok {
		return "", fmt.Errorf("you can't go %s from here", dir)
	}

	// Leave current room
	p.Room.RemovePlayer(p.Name)

	// Join newly entered room
	targetRoom.AddPlayer(p)
	p.Room = targetRoom

	return fmt.Sprintf("---\nYou move %s into %s.\n", dir, targetRoom.Name), nil
}
