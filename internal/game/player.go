package game

import (
	"fmt"
	"log"
	"net"
	"strings"
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
	Name      string
	Conn      net.Conn
	Room      *Room
	Stats     Stats
	Inventory []string
	Equipped  []string
	//Tracks which stats have been set
	AssignedStats map[string]bool
}

func (p *Player) ShowStats() string {
	stats := p.Stats

	statsBlock := "\n<!!--------------------------------!!>"
	statsBlock += "\n+-- Stats for %s"
	statsBlock += "\n<!!--------------------------------!!>"
	statsBlock += "\n+-- STR: %d   DEX: %d   CON: %d"
	statsBlock += "\n+-- INT: %d   WIS: %d   CHA: %d"
	statsBlock += "\n<!!--------------------------------!!>"

	return fmt.Sprintf(statsBlock, p.Name, stats.STR, stats.DEX, stats.CON, stats.INT, stats.WIS, stats.CHA)
}

func (p *Player) AssignStat(stat string) (string, error) {
	stat = strings.ToUpper(stat)

	log.Printf("Assigned stat: %s", stat)
	log.Printf("Assigned stats: %v", p.AssignedStats)
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
