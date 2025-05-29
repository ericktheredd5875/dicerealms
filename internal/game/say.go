package game

import "fmt"

func (p *Player) Say(text string) string {
	if text == "" {
		return ""
	}

	// Send to the person who spoke
	uSpoke := fmt.Sprintf(`You say, "%s"`+"\n", text)
	uSpoke = Colorize(uSpoke, Green)
	p.Conn.Write([]byte(uSpoke))

	// Broadcast to the room
	msg := fmt.Sprintf(`%s says, "%s"`, p.Name, text)
	p.Room.Broadcast(Colorize(msg, Green), p.Name)

	if p.Room.ActiveScene != nil {
		p.Room.ActiveScene.LogEntry(msg)
	}

	return msg
}
