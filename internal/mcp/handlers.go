package mcp

import (
	"fmt"
	"log"
	"net"

	"github.com/ericktheredd5875/dicerealms/internal/game"
	"github.com/ericktheredd5875/dicerealms/internal/session"
)

func HandleMCPRegister(args map[string]string, player *game.Player, conn net.Conn) {
	name := args["name"]
	if name == "" {
		conn.Write([]byte(`Usage: #$#mcp-register name="<character_name>"`))
		return
	}

	room := player.Room
	err := player.RegisterPlayer(name)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("Registration failed: %s", err)))
		log.Printf("Registration failed: %s", err)
		return
	}

	conn.Write([]byte(fmt.Sprintf("Welcome, %s! Thank you for registering your character.", player.Name)))
	game.JoinRoom(player, room, conn)
}

func HandleMCPLogin(args map[string]string, player *game.Player, conn net.Conn) {
	name := args["name"]
	if name == "" {
		conn.Write([]byte(`Usage: #$#mcp-login name="<character_name>"`))
		return
	}

	loginPlayer, err := game.HandleLogin(name)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("Login failed: %s", err)))
		log.Printf("Login failed: %s", err)
		return
	}

	room := player.Room
	player = loginPlayer
	player.Room = room
	game.JoinRoom(player, room, conn)

	// player.Conn = conn

	log.Printf("Player: %+v", player)
	conn.Write([]byte(fmt.Sprintf("Welcome back, %s!\n", player.Name)))
}

func HandleMCPExit(args map[string]string, player *game.Player, s *session.Session) {
	// if player != nil {
	// 	player.Save()
	// 	player.Room.RemovePlayer(player.Name)
	// }

	// conn.Write([]byte("Farewell, travler.... Come Again"))
	// conn.Close()
	// log.Printf("Player: %+v", player)
	// s.Send("Farewell, travler.... Come Again")
	s.Close()
}
