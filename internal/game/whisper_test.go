package game

import (
	"net"
	"strings"
	"testing"
	"time"
)

type whisperBufferConn struct {
	data strings.Builder
}

func (b *whisperBufferConn) Write(p []byte) (n int, err error) {
	return b.data.Write(p)
}

func (b *whisperBufferConn) Read(p []byte) (int, error)         { return 0, nil }
func (b *whisperBufferConn) Close() error                       { return nil }
func (b *whisperBufferConn) LocalAddr() net.Addr                { return nil }
func (b *whisperBufferConn) RemoteAddr() net.Addr               { return nil }
func (b *whisperBufferConn) SetDeadline(t time.Time) error      { return nil }
func (b *whisperBufferConn) SetReadDeadline(t time.Time) error  { return nil }
func (b *whisperBufferConn) SetWriteDeadline(t time.Time) error { return nil }

func TestWhisper(t *testing.T) {
	room := NewRoom("Hall", "Echoing stone chamber.")
	aConn := &whisperBufferConn{}
	bConn := &whisperBufferConn{}

	alice := &Player{Name: "Alice", Conn: aConn}
	bob := &Player{Name: "Bob", Conn: bConn}

	room.AddPlayer(alice)
	room.AddPlayer(bob)

	alice.Room = room
	bob.Room = room

	err := alice.Whisper(bob.Name, "This is a test.")
	if err != nil {
		t.Errorf("Alice did not get correct whisper: %q", err)
	}

	if !strings.Contains(aConn.data.String(), `You whisper to Bob: "This is a test."`) {
		t.Errorf("Alice did not get correct echo: %q", aConn.data.String())
	}

	if !strings.Contains(bConn.data.String(), `Alice whispers: "This is a test."`) {
		t.Errorf("Bob did not get correct whisper: %q", bConn.data.String())
	}

	if err != nil {
		t.Errorf("Returned message incorrect: %q", err)
	}
}
