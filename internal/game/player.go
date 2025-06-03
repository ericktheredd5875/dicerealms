package game

import (
	"fmt"
	"net"
	"strings"

	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
	"gorm.io/gorm"
)

type Stats struct {
	STR int
	DEX int
	CON int
	INT int
	WIS int
	CHA int
}

type Player struct {
	ID        uint
	PublicID  string
	Name      string
	Conn      net.Conn
	Room      *Room
	Stats     Stats
	Inventory []string
	Equipped  []string
	Gold      int
	Level     int
	XP        int
	//Tracks which stats have been set
	AssignedStats map[string]bool
}

func PlayerPrompt(playerName string, roomName string) string {

	prompt := Colorize("%s@%s +>>", Bold+Cyan)
	prompt = fmt.Sprintf("\n"+prompt, playerName, roomName)
	return prompt
}

func (p *Player) ShowStats() string {
	stats := p.Stats

	statsBlock := "<!!--------------------------------!!>"
	statsBlock += "\n+-- Stats for %s"
	statsBlock += "\n<!!--------------------------------!!>"
	statsBlock += "\n+-- STR: %d   DEX: %d   CON: %d"
	statsBlock += "\n+-- INT: %d   WIS: %d   CHA: %d"
	statsBlock += "\n<!!--------------------------------!!>"

	return fmt.Sprintf(statsBlock, p.Name, stats.STR, stats.DEX, stats.CON, stats.INT, stats.WIS, stats.CHA)
}

func (p *Player) AssignStat(stat string) (string, error) {
	stat = strings.ToUpper(stat)

	// log.Printf("Assigned stat: %s", stat)
	// log.Printf("Assigned stats: %v", p.AssignedStats)
	if p.AssignedStats[stat] {
		return "", fmt.Errorf("stat already assigned: %s", stat)
	}

	val, rolls := roll4d6DropLowest()

	switch stat {
	case "STR":
		p.Stats.STR = val
	case "DEX":
		p.Stats.DEX = val
	case "CON":
		p.Stats.CON = val
	case "INT":
		p.Stats.INT = val
	case "WIS":
		p.Stats.WIS = val
	case "CHA":
		p.Stats.CHA = val
	default:
		return "", fmt.Errorf("invalid stat: %s", stat)
	}

	p.AssignedStats[stat] = true
	return fmt.Sprintf("Rolled 4d6 -> %v. Set %s = %d", rolls, stat, val), nil
}

func (p *Player) AutoGenStats() (string, error) {

	for _, s := range validStats {
		if p.AssignedStats[s] {
			return "", fmt.Errorf("some stats have already been assigned -- cannot auto-generate")
		}
	}

	results := []string{}
	for _, stat := range validStats {
		msg, _ := p.AssignStat(stat)
		results = append(results, msg)
	}

	return strings.Join(results, "\n"), nil
}

func (p *Player) AddItem(item string) {
	p.Inventory = append(p.Inventory, item)

}

func (p *Player) RemoveItem(item string) bool {
	for i, v := range p.Inventory {
		if v == item {
			p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
			return true
		}
	}

	return false
}

func (p *Player) InventoryList() string {
	if len(p.Inventory) == 0 {
		return "Your inventory is empty."
	}

	list := "Inventory: \n"
	for _, item := range p.Inventory {
		list += fmt.Sprintf("+-- %s\n", item)
	}

	return list
}

func (p *Player) Save() error {
	model := db.PlayerModel{
		Model:     gorm.Model{ID: p.ID},
		PublicID:  p.PublicID,
		Name:      p.Name,
		STR:       p.Stats.STR,
		DEX:       p.Stats.DEX,
		CON:       p.Stats.CON,
		INT:       p.Stats.INT,
		WIS:       p.Stats.WIS,
		CHA:       p.Stats.CHA,
		Inventory: p.Inventory,
		Equipped:  p.Equipped,
		Gold:      p.Gold,
		Level:     p.Level,
		XP:        p.XP,
		// RoomID: p.Room.ID,
	}

	return db.DB.Save(&model).Error
}

func HandleLogin(name string) (*Player, error) {
	var model db.PlayerModel
	if err := db.DB.Where("name = ?", name).First(&model).Error; err != nil {
		return nil, err
	}

	player := &Player{
		ID:       model.ID,
		PublicID: model.PublicID,
		Name:     model.Name,
		Stats: Stats{
			STR: model.STR,
			DEX: model.DEX,
			CON: model.CON,
			INT: model.INT,
			WIS: model.WIS,
			CHA: model.CHA,
		},
		Inventory:     model.Inventory,
		Equipped:      model.Equipped,
		Gold:          model.Gold,
		Level:         model.Level,
		XP:            model.XP,
		AssignedStats: map[string]bool{},
	}

	for _, stat := range validStats {
		switch stat {
		case "STR":
			if player.Stats.STR > 0 {
				player.AssignedStats[stat] = true
			}
		case "DEX":
			if player.Stats.DEX > 0 {
				player.AssignedStats[stat] = true
			}
		case "CON":
			if player.Stats.CON > 0 {
				player.AssignedStats[stat] = true
			}
		case "INT":
			if player.Stats.INT > 0 {
				player.AssignedStats[stat] = true
			}
		case "WIS":
			if player.Stats.WIS > 0 {
				player.AssignedStats[stat] = true
			}
		case "CHA":
			if player.Stats.CHA > 0 {
				player.AssignedStats[stat] = true
			}
		}
	}

	return player, nil
}

func (p *Player) RegisterPlayer(name string) error {
	pubID, _ := utils.GenerateKHash(name+":DiceRealms:Telnet:4000", "")
	model := db.PlayerModel{
		Name:     name,
		PublicID: pubID,
	}

	if err := db.DB.Create(&model).Error; err != nil {
		return err
	}

	p.ID = model.ID
	p.PublicID = model.PublicID
	p.Name = model.Name

	return nil
	// return &Player{
	// 	ID:       model.ID,
	// 	PublicID: model.PublicID,
	// 	Name:     model.Name,
	// 	// Stats: Stats{
	// 	// 	STR: model.STR,
	// 	// 	DEX: model.DEX,
	// 	// 	CON: model.CON,
	// 	// 	INT: model.INT,
	// 	// 	WIS: model.WIS,
	// 	// 	CHA: model.CHA,
	// 	// },
	// 	// Inventory:     model.Inventory,
	// 	// Equipped:      model.Equipped,
	// 	// Gold:          model.Gold,
	// 	// Level:         model.Level,
	// 	// XP:            model.XP,
	// 	AssignedStats: map[string]bool{},
	// }, nil
}

func JoinRoom(player *Player, room *Room, conn net.Conn) {
	// player.Room = room
	player.Conn = conn
	room.AddPlayer(player)
	conn.Write([]byte(fmt.Sprintf("Welcome to %s!", room.Name)))

}
