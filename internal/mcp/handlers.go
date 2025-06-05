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

func HandleRegister(args map[string]string, p *game.Player, s *session.Session) {
	name := args["name"]
	if name == "" {
		s.Send(utils.Colorize("Usage: #$#mcp-register name=\"<character_name>\"\n", utils.BrGreen))
		return
	}

	room := p.Room
	err := p.RegisterPlayer(name)
	if err != nil {
		s.Send(utils.ColorizeError(fmt.Sprintf("Registration failed: %s", err)))
		log.Printf("Registration failed: %s", err)
		return
	}

	s.Send(fmt.Sprintf("Welcome, %s! Thank you for registering your character.", p.Name))
	game.JoinRoom(p, room, s.Conn)
}

func HandleLogin(args map[string]string, p *game.Player, s *session.Session) {
	name := args["name"]
	if name == "" {
		s.Send(utils.Colorize("Usage: #$#mcp-login name=\"<character_name>\"\n", utils.BrGreen))
		return
	}

	loginPlayer, err := game.HandleLogin(name)
	if err != nil {
		s.Send(utils.ColorizeError(fmt.Sprintf("Login failed: %s", err)))
		log.Printf("Login failed: %s", err)
		return
	}

	room := p.Room
	p = loginPlayer
	p.Room = room
	game.JoinRoom(p, room, s.Conn)

	// p.Conn = conn

	log.Printf("Player: %+v", p)
	s.Send(fmt.Sprintf("Welcome back, %s!\n", p.Name))
}

func HandleGo(args map[string]string, p *game.Player, s *session.Session) {
	dir := strings.ToLower(strings.TrimSpace(args["dir"]))
	if dir == "" {

		msg := utils.Colorize("Usage: #$#mcp-go dir=\"north\"\n", utils.BrGreen)
		s.Send(msg)
		return
	}

	curRoom := p.Room
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
	curRoom.RemovePlayer(p.Name)

	// Update Player room
	p.Room = nextRoom

	// Save and Persist to DB
	if p.Model != nil {
		p.Model.RoomID = nextRoom.ID
		db.DB.Save(p.Model)
	}

	s.Send(fmt.Sprintf("You move %s into %s.\n", dir, nextRoom.Name))
	s.Send(nextRoom.Desc)

}

func HandlePickup(args map[string]string, p *game.Player, s *session.Session) {
	// name := strings.ToLower(strings.TrimSpace(args["name"]))
	name := strings.TrimSpace(args["name"])
	if name == "" {
		s.Send(utils.Colorize("Usage: #$#mcp-pickup name=\"<item_name>\"\n", utils.BrGreen))
		return
	}

	item := game.GetItemByName(name)
	if item == nil {
		s.Send(utils.Colorize(fmt.Sprintf("[x] Item not found: %s\n", name), utils.BrRed))
		return
	}

	// Add to Player Inventory
	p.Inventory = append(p.Inventory, item.Name)
	item.RoomFoundID = 0
	db.DB.Save(item)
	p.Save()

	s.Send(utils.Colorize(fmt.Sprintf("[+] You pick up %s.\n", name), utils.BrGreen))
}

func HandleDrop(args map[string]string, p *game.Player, s *session.Session) {
	// name := strings.ToLower(strings.TrimSpace(args["name"]))
	name := strings.TrimSpace(args["name"])
	if name == "" {
		s.Send(utils.Colorize("Usage: #$#mcp-drop name=\"<item_name>\"\n", utils.BrGreen))
		return
	}

	found := false
	newInventory := []string{}
	for _, i := range p.Inventory {
		if i == name && !found {
			found = true
			continue
		}

		newInventory = append(newInventory, i)
	}

	if !found {
		s.Send(utils.Colorize(fmt.Sprintf("[x] Item not found in inventory: %s\n", name), utils.BrRed))
		return
	}

	p.Inventory = newInventory
	p.Save()

	item := game.GetItemByName(name)
	if item != nil {
		item.RoomFoundID = uint(p.RoomID)
		db.DB.Save(item)
	}

	s.Send(utils.Colorize(fmt.Sprintf("[+] You drop %s.\n", name), utils.BrGreen))
}

func HandleExamine(args map[string]string, p *game.Player, s *session.Session) {
	name := strings.TrimSpace(args["name"])
	if name == "" {
		s.Send(utils.Colorize("Usage: #$#mcp-examine name=\"<item_name>\"\n", utils.BrGreen))
		return
	}

	// Check Inventory
	for _, i := range p.Inventory {
		if strings.EqualFold(i, name) {
			if i := game.GetItemByName(i); i != nil {
				s.Send(fmt.Sprintf("%s; %s", i.Name, i.Description))
				return
			}
		}
	}

	roomItem := &db.ItemModel{}
	for _, i := range game.GetAllItemsInRoom(uint(p.RoomID)) {
		if strings.EqualFold(i.Name, name) {
			roomItem = i
			break
		}
	}

	if roomItem != nil {
		s.Send(fmt.Sprintf("%s; %s", roomItem.Name, roomItem.Description))
	} else {
		s.Send("You don't see that item here or in your inventory.")
	}
}

func HandleInventory(p *game.Player, s *session.Session) {
	if len(p.Inventory) == 0 {
		s.Send("Your inventory is empty.")
		return
	}

	list := "You are Carrying: \n"
	for _, item := range p.Inventory {
		if i := game.GetItemByName(item); i != nil {
			list += fmt.Sprintf("+-- %s [Effect: %s]\n", i.Name, i.Effect)
		}
	}

	msg := utils.ColorizeSuccess(list + "\n")
	s.Send(msg)

	s.Send(fmt.Sprintf("Gold: %d", p.Gold))
}

func HandleExit(args map[string]string, p *game.Player, s *session.Session) {

	s.Close()
}
