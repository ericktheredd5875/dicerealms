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

	whisperMsg := Colorize(fmt.Sprintf(`You whisper to %s: "%s"`, to, msg), Cyan)
	p.Conn.Write([]byte(whisperMsg + "\n"))

	targetWhisperMsg := Colorize(fmt.Sprintf(`%s whispers: "%s"`, p.Name, msg), Cyan)
	target.Conn.Write([]byte("\n" + targetWhisperMsg + "\n"))
	targetPrompt := PlayerPrompt(target.Name, target.Room.Name)
	target.Conn.Write([]byte(targetPrompt))

	return nil
}
