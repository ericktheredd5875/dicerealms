package mcp

import (
	"fmt"
	"log"
	"strings"

	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/ericktheredd5875/dicerealms/internal/game"
	"github.com/ericktheredd5875/dicerealms/internal/session"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

func HandleMCPRegister(args map[string]string, player *game.Player, s *session.Session) {
	name := args["name"]
	if name == "" {
		s.Send(`Usage: #$#mcp-register name="<character_name>"`)
		return
	}

	room := player.Room
	err := player.RegisterPlayer(name)
	if err != nil {
		s.Send(utils.ColorizeError(fmt.Sprintf("Registration failed: %s", err)))
		log.Printf("Registration failed: %s", err)
		return
	}

	s.Send(fmt.Sprintf("Welcome, %s! Thank you for registering your character.", player.Name))
	game.JoinRoom(player, room, s.Conn)
}

func HandleMCPLogin(args map[string]string, player *game.Player, s *session.Session) {
	name := args["name"]
	if name == "" {
		s.Send(`Usage: #$#mcp-login name="<character_name>"`)
		return
	}

	loginPlayer, err := game.HandleLogin(name)
	if err != nil {
		s.Send(utils.ColorizeError(fmt.Sprintf("Login failed: %s", err)))
		log.Printf("Login failed: %s", err)
		return
	}

	room := player.Room
	player = loginPlayer
	player.Room = room
	game.JoinRoom(player, room, s.Conn)

	// player.Conn = conn

	log.Printf("Player: %+v", player)
	s.Send(fmt.Sprintf("Welcome back, %s!\n", player.Name))
}

func HandleMCPGo(args map[string]string, player *game.Player, s *session.Session) {
	dir := strings.ToLower(strings.TrimSpace(args["dir"]))
	if dir == "" {

		msg := utils.Colorize("Usage: #$#mcp-go dir=\"north\"\n", utils.BrGreen)
		s.Send(msg)
		return
	}

	curRoom := player.Room
	if curRoom == nil {
		msg := utils.Colorize("[x] You're lost in the void. No Room found.\n", utils.BrRed)
		s.Send(msg)
		return
	}

	curRoom.Mu.Lock()
	nextRoom, ok := curRoom.Exits[dir]
	curRoom.Mu.Unlock()

	if !ok {
		msg := utils.Colorize(fmt.Sprintf("[x] No Exit in direction: %s.\n", dir), utils.BrRed)
		s.Send(msg)
		return
	}

	// Remove Player from Current Room
	curRoom.RemovePlayer(player.Name)

	// Update Player room
	player.Room = nextRoom

	// Save and Persist to DB
	if player.Model != nil {
		player.Model.RoomID = nextRoom.ID
		db.DB.Save(player.Model)
	}

	s.Send(fmt.Sprintf("You move %s into %s.\n", dir, nextRoom.Name))
	s.Send(nextRoom.Desc)

}

func HandleMCPExit(args map[string]string, player *game.Player, s *session.Session) {

	s.Close()
}
