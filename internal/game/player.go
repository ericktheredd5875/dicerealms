package game

import "net"

type Player struct {
	Name string
	Conn net.Conn
	Room *Room
}
