package game

import (
	"fmt"
)

func (p *Player) Whisper(to string, msg string) error {
	p.Room.mu.Lock()
	target, ok := p.Room.Players[to]
	p.Room.mu.Unlock()

	if !ok {
		return fmt.Errorf("%s is not in the room", to)
	}

	p.Conn.Write([]byte(fmt.Sprintf(`You whisper to %s: "%s"`+"\n", to, msg)))
	target.Conn.Write([]byte(fmt.Sprintf(`%s whispers: "%s"`+"\n", p.Name, msg)))
	target.Conn.Write([]byte("+>> "))

	return nil
}
