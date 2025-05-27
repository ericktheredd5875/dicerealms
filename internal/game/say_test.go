package game

import (
	"net"
	"strings"
	"testing"
	"time"
)

type sayBufferConn struct {
	data strings.Builder
}

func (b *sayBufferConn) Write(p []byte) (n int, err error) {
	return b.data.Write(p)
}

func (b *sayBufferConn) Read(p []byte) (int, error)         { return 0, nil }
func (b *sayBufferConn) Close() error                       { return nil }
func (b *sayBufferConn) LocalAddr() net.Addr                { return nil }
func (b *sayBufferConn) RemoteAddr() net.Addr               { return nil }
func (b *sayBufferConn) SetDeadline(t time.Time) error      { return nil }
func (b *sayBufferConn) SetReadDeadline(t time.Time) error  { return nil }
func (b *sayBufferConn) SetWriteDeadline(t time.Time) error { return nil }

func TestSay(t *testing.T) {
	room := NewRoom("Hall", "Echoing stone chamber.")
	aConn := &sayBufferConn{}
	bConn := &sayBufferConn{}

	alice := &Player{Name: "Alice", Conn: aConn}
	bob := &Player{Name: "Bob", Conn: bConn}

	room.AddPlayer(alice)
	room.AddPlayer(bob)

	alice.Room = room
	bob.Room = room

	msg := alice.Say("This is a test.")

	if !strings.Contains(aConn.data.String(), `You say, "This is a test."`) {
		t.Errorf("Alice did not get correct echo: %q", aConn.data.String())
	}

	if !strings.Contains(bConn.data.String(), `Alice says, "This is a test."`) {
		t.Errorf("Bob did not get correct broadcast: %q", bConn.data.String())
	}

	if msg != `Alice says, "This is a test."` {
		t.Errorf("Returned message incorrect: %q", msg)
	}
}
