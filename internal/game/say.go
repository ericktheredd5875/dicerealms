package game

import "fmt"

func (p *Player) Say(text string) string {
	if text == "" {
		return ""
	}

	// Send to the person who spoke
	p.Conn.Write([]byte(fmt.Sprintf(`You say, "%s"`+"\n", text)))

	// Broadcast to the room
	msg := fmt.Sprintf(`%s says, "%s"`, p.Name, text)
	p.Room.Broadcast(msg, p.Name)

	return msg
}
