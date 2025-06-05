package netiface

import "net"

type GameConn interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	ReadLine() (string, error)
	RemoteAddr() net.Addr
}
