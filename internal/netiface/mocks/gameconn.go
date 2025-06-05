package mocks

import (
	"bufio"
	"net"
	"strings"

	"github.com/stretchr/testify/mock"
)

// GameConn is a mock of netiface.GameConn
type GameConn struct {
	mock.Mock
	Conn net.Conn
}

func (m *GameConn) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *GameConn) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *GameConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *GameConn) RemoteAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *GameConn) ReadLine() (string, error) {
	line, err := bufio.NewReader(m.Conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
