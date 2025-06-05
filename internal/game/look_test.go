package game

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/ericktheredd5875/dicerealms/internal/netiface"
)

type dummyConn struct{}

func (d *dummyConn) Read(b []byte) (int, error)         { return 0, nil }
func (d *dummyConn) Write(b []byte) (int, error)        { return len(b), nil }
func (d *dummyConn) Close() error                       { return nil }
func (d *dummyConn) LocalAddr() net.Addr                { return nil }
func (d *dummyConn) RemoteAddr() net.Addr               { return nil }
func (d *dummyConn) SetDeadline(t time.Time) error      { return nil }
func (d *dummyConn) SetReadDeadline(t time.Time) error  { return nil }
func (d *dummyConn) SetWriteDeadline(t time.Time) error { return nil }

func TestLook(t *testing.T) {
	room := NewRoom(
		"Hall of Echoes",
		"A long, shadowy corridor with whispers in the dark.")
	otherRoom := NewRoom(
		"Outside",
		"A bright field under a blue sky.")
	room.Exits["east"] = otherRoom

	alice := &Player{Name: "Alice", Conn: &netiface.TelnetConn{Conn: &dummyConn{}}}
	bob := &Player{Name: "Bob", Conn: &netiface.TelnetConn{Conn: &dummyConn{}}}

	room.AddPlayer(alice)
	room.AddPlayer(bob)

	alice.Room = room
	look := alice.Look()

	// Verify the output
	if !strings.Contains(look, "Hall of Echoes") {
		t.Errorf("Room name missing from look output")
	}
	if !strings.Contains(look, "A long, shadowy corridor") {
		t.Errorf("Room description missing")
	}
	if !strings.Contains(look, "- Bob") {
		t.Errorf("Other player not listed")
	}
	if !strings.Contains(look, "- east") {
		t.Errorf("Exit not listed")
	}
}
