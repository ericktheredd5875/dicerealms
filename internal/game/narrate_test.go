package game

import (
	"net"
	"strings"
	"testing"
	"time"
)

type narrateConn struct {
	strings.Builder
}

func (l *narrateConn) Write(b []byte) (int, error)        { return l.Builder.Write(b) }
func (l *narrateConn) Read(b []byte) (int, error)         { return 0, nil }
func (l *narrateConn) Close() error                       { return nil }
func (l *narrateConn) LocalAddr() net.Addr                { return nil }
func (l *narrateConn) RemoteAddr() net.Addr               { return nil }
func (l *narrateConn) SetDeadline(t time.Time) error      { return nil }
func (l *narrateConn) SetReadDeadline(t time.Time) error  { return nil }
func (l *narrateConn) SetWriteDeadline(t time.Time) error { return nil }

func TestNarrate(t *testing.T) {
	room := NewRoom("Test Room", "Just a test.")
	aConn := &narrateConn{}
	bConn := &narrateConn{}

	alice := &Player{Name: "Alice", Room: room, Conn: aConn}
	bob := &Player{Name: "Bob", Room: room, Conn: bConn}

	room.AddPlayer(alice)
	room.AddPlayer(bob)

	alice.Narrate("A shadow moves across the room.", alice.Name)

	if !strings.Contains(bConn.String(), "<narritive> A shadow moves across the room.") {
		t.Errorf("Bob did not receive the message: %s", bConn.String())
	}

	if strings.Contains(aConn.String(), "Alice whispers") {
		t.Errorf("Narration should be system-style, not character-attributed. %s", aConn.String())
	}
}
