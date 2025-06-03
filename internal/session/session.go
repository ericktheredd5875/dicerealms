package session

import (
	"bufio"
	"log"
	"net"

	"github.com/ericktheredd5875/dicerealms/internal/game"
)

type Session struct {
	Conn   net.Conn
	Writer *bufio.Writer
	Player *game.Player
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		Conn:   conn,
		Writer: bufio.NewWriter(conn),
	}
}

func (s *Session) Send(msg string) {
	s.Writer.WriteString(msg)
	s.Writer.Flush()
}

func (s *Session) Close() {
	// log.Printf("Closing session for player: %s", s.Player.Name)
	log.Printf("Player: %+v", s.Player)
	if s.Player != nil {
		s.Send(game.ColorizeInfo("!!Farewell, travler.... Come Again!!"))
		s.Player.Save()
		s.Player.Room.RemovePlayer(s.Player.Name)
	}

	s.Conn.Close()
}
