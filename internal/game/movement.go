package game

import (
	"fmt"

	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

func (p *Player) Move(dir string) (string, error) {
	p.Room.Mu.Lock()
	targetRoom, ok := p.Room.Exits[dir]
	p.Room.Mu.Unlock()

	if !ok {
		return "", fmt.Errorf("you can't go %s from here", dir)
	}

	// Leave current room
	p.Room.RemovePlayer(p.Name)

	// Join newly entered room
	targetRoom.AddPlayer(p)
	p.Room = targetRoom

	msg := fmt.Sprintf("---\nYou move %s into %s.\n", dir, targetRoom.Name)
	return utils.Colorize(msg, utils.Gray), nil
}
