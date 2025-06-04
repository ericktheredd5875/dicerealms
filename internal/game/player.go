package game

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
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
	Model         *db.PlayerModel
	LastActiveAt  time.Time
	RoomID        int
}

func PlayerPrompt(playerName string, roomName string) string {

	prompt := utils.Colorize("%s@%s +>>", utils.Bold+utils.Cyan)
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
	// model := db.PlayerModel{
	// 	Model:     gorm.Model{ID: p.ID},
	// 	PublicID:  p.PublicID,
	// 	Name:      p.Name,
	// 	STR:       p.Stats.STR,
	// 	DEX:       p.Stats.DEX,
	// 	CON:       p.Stats.CON,
	// 	INT:       p.Stats.INT,
	// 	WIS:       p.Stats.WIS,
	// 	CHA:       p.Stats.CHA,
	// 	Inventory: p.Inventory,
	// 	Equipped:  p.Equipped,
	// 	Gold:      p.Gold,
	// 	Level:     p.Level,
	// 	XP:        p.XP,
	// 	RoomID:    p.Room.ID,
	// }

	// p.Conn.Write([]byte(fmt.Sprintf("Saving player %s...", p.Name)))
	log.Printf("Saving player %s...", p.Name)
	UpdateModelFromPlayer(p)
	log.Printf("Player model: %+v", p.Model)

	return db.DB.Save(p.Model).Error
}

func HandleLogin(name string) (*Player, error) {
	var model db.PlayerModel
	if err := db.DB.Where("name = ?", name).First(&model).Error; err != nil {
		return nil, err
	}

	player := ToPlayer(&model)

	return player, nil
}

func (p *Player) RegisterPlayer(name string) error {
	pubID, _ := utils.GenerateKHash(name+":DiceRealms:Telnet:4000", "")
	model := db.PlayerModel{
		Name:     name,
		PublicID: pubID,
	}

	log.Printf("Registering player %s...", name)
	log.Printf("Model: %+v", model)

	if err := db.DB.Create(&model).Error; err != nil {
		return err
	}

	p.ID = model.ID
	p.PublicID = model.PublicID
	p.Name = model.Name
	p.Model = &model

	return nil
}

func JoinRoom(player *Player, room *Room, conn net.Conn) {

	player.Room = room
	player.Conn = conn
	room.AddPlayer(player)
	conn.Write([]byte(fmt.Sprintf("Welcome to %s!", room.Name)))

}

func ToPlayer(model *db.PlayerModel) *Player {
	stats := map[string]int{
		"STR": int(model.STR),
		"DEX": int(model.DEX),
		"CON": int(model.CON),
		"INT": int(model.INT),
		"WIS": int(model.WIS),
		"CHA": int(model.CHA),
	}

	assignedStats := map[string]bool{}
	for _, stat := range validStats {
		if stats[stat] > 0 {
			assignedStats[stat] = true
		}
	}

	return &Player{
		ID:       model.ID,
		PublicID: model.PublicID,
		Name:     model.Name,
		Stats: Stats{
			STR: stats["STR"],
			DEX: stats["DEX"],
			CON: stats["CON"],
			INT: stats["INT"],
			WIS: stats["WIS"],
		},
		Inventory:     model.Inventory,
		Equipped:      model.Equipped,
		Gold:          model.Gold,
		Level:         model.Level,
		XP:            model.XP,
		AssignedStats: assignedStats,
		RoomID:        model.LastRoomID,
		Model:         model,
	}
}

func UpdateModelFromPlayer(player *Player) {

	if player.Model == nil {
		return
	}

	log.Printf("Updating model from player %s...", player.Name)

	player.Model.STR = int(player.Stats.STR)
	player.Model.DEX = int(player.Stats.DEX)
	player.Model.CON = int(player.Stats.CON)
	player.Model.INT = int(player.Stats.INT)
	player.Model.WIS = int(player.Stats.WIS)
	player.Model.CHA = int(player.Stats.CHA)

	player.Model.Inventory = player.Inventory
	player.Model.Equipped = player.Equipped
	player.Model.Gold = player.Gold
	player.Model.Level = player.Level
	player.Model.XP = player.XP

	if player.Room != nil {
		player.Model.LastRoomID = player.Room.ID
	}
}
