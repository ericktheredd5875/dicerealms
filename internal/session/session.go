package session

import (
	"bufio"

	"github.com/ericktheredd5875/dicerealms/internal/game"
	"github.com/ericktheredd5875/dicerealms/internal/netiface"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

type Session struct {
	Conn   netiface.GameConn
	Writer *bufio.Writer
	Player *game.Player
}

func NewSession(conn netiface.GameConn) *Session {
	return &Session{
		Conn:   conn,
		Writer: bufio.NewWriter(conn),
	}
}

func (s *Session) Send(msg string) {

	// s.Conn.Write([]byte(msg))
	s.Writer.WriteString(msg)
	s.Writer.Flush()
}

func (s *Session) Close() {
	// log.Printf("SessionClose-Player: %+v", s.Player)
	if s.Player != nil {
		s.Send(utils.ColorizeInfo("!!Farewell, travler.... Come Again!!") + "\n")
		s.Player.Save()
		s.Player.Room.RemovePlayer(s.Player.Name)
	}

	s.Conn.Close()
}
