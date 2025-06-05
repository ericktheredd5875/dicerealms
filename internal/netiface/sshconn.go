package netiface

import (
	"net"

	"github.com/gliderlabs/ssh"
)

type SSHConn struct {
	Session ssh.Session
}

func (c *SSHConn) Read(p []byte) (int, error)  { return c.Session.Read(p) }
func (c *SSHConn) Write(p []byte) (int, error) { return c.Session.Write(p) }
func (c *SSHConn) Close() error                { return c.Session.Close() }
func (c *SSHConn) RemoteAddr() net.Addr        { return c.Session.RemoteAddr() }

func (c *SSHConn) ReadLine() (string, error) {
	var buf []byte
	tmp := make([]byte, 1)

	for {
		n, err := c.Session.Read(tmp)
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}

		ch := tmp[0]
		if ch == '\r' || ch == '\n' {
			c.Session.Write([]byte("\r\n"))
			break
		}
		if ch == 127 || ch == 8 {
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				c.Session.Write([]byte("\b \b"))
			}
			continue
		}

		buf = append(buf, ch)
		c.Session.Write([]byte{ch})
	}

	return string(buf), nil
}
