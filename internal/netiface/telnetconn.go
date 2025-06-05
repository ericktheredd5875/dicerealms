package netiface

import (
	"bufio"
	"net"
	"strings"
)

type TelnetConn struct {
	Conn net.Conn
}

func (c *TelnetConn) Read(p []byte) (int, error)  { return c.Conn.Read(p) }
func (c *TelnetConn) Write(p []byte) (int, error) { return c.Conn.Write(p) }
func (c *TelnetConn) Close() error                { return c.Conn.Close() }

func (c *TelnetConn) RemoteAddr() net.Addr { return c.Conn.RemoteAddr() }

func (c *TelnetConn) ReadLine() (string, error) {
	line, err := bufio.NewReader(c.Conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
