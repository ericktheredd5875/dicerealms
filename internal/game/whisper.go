package game

import (
	"fmt"

	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

func (p *Player) Whisper(to string, msg string) error {
	p.Room.Mu.Lock()
	target, ok := p.Room.Players[to]
	p.Room.Mu.Unlock()

	if !ok {
		return fmt.Errorf("%s is not in the room", to)
	}

	whisperMsg := utils.Colorize(fmt.Sprintf(`You whisper to %s: "%s"`, to, msg), utils.Cyan)
	p.Conn.Write([]byte(whisperMsg + "\n"))

	targetWhisperMsg := utils.Colorize(fmt.Sprintf(`%s whispers: "%s"`, p.Name, msg), utils.Cyan)
	target.Conn.Write([]byte("\n" + targetWhisperMsg + "\n"))
	targetPrompt := PlayerPrompt(target.Name, target.Room.Name)
	target.Conn.Write([]byte(targetPrompt))

	return nil
}
