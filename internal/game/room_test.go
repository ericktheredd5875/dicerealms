package game

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"
)

type mockConn struct {
	bytes.Buffer
}

func (m *mockConn) Write(b []byte) (int, error) {
	return m.Buffer.Write(b)
}

func (m *mockConn) Read(b []byte) (int, error)         { return 0, nil }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestRoomBroadcast(t *testing.T) {
	room := NewRoom("Test Room")

	aliceConn := &mockConn{}
	bobConn := &mockConn{}

	alice := &Player{Name: "Alice", Conn: aliceConn}
	bob := &Player{Name: "Bob", Conn: bobConn}

	room.AddPlayer(alice)
	room.AddPlayer(bob)

	room.Broadcast("The dragon roars!", "Alice")

	got := bobConn.String()
	want := "The dragon roars!\n"

	if got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}

	if aliceConn.String() != "" && aliceConn.String() == want {
		log.Printf("Alice received message: %q", aliceConn.String())
		t.Errorf("Alice should not receive her own message")
	}
}
