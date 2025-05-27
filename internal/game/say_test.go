package game

import (
	"net"
	"strings"
	"testing"
	"time"
)

type bufferConn struct {
	data strings.Builder
}

func (b *bufferConn) Write(p []byte) (n int, err error) {
	return b.data.Write(p)
}

func (b *bufferConn) Read(p []byte) (int, error)         { return 0, nil }
func (b *bufferConn) Close() error                       { return nil }
func (b *bufferConn) LocalAddr() net.Addr                { return nil }
func (b *bufferConn) RemoteAddr() net.Addr               { return nil }
func (b *bufferConn) SetDeadline(t time.Time) error      { return nil }
func (b *bufferConn) SetReadDeadline(t time.Time) error  { return nil }
func (b *bufferConn) SetWriteDeadline(t time.Time) error { return nil }

func TestSay(t *testing.T) {
	room := NewRoom("Hall", "Echoing stone chamber.")
	aConn := &bufferConn{}
	bConn := &bufferConn{}

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
