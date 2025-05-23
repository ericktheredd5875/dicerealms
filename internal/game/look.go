package game

import (
	"fmt"
	"log"
	"strings"
)

func (p *Player) Look() string {
	r := p.Room
	r.mu.Lock()
	defer r.mu.Unlock()

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("---\n|- %s\n %s\n---\n", r.Name, r.Desc))

	// Players
	builder.WriteString("|- Players here:\n")
	foundOther := false
	for name := range r.Players {
		if name != p.Name {
			builder.WriteString("|-- " + name + "\n")
			foundOther = true
		}
	}

	if !foundOther {
		builder.WriteString("|-- (You are alone here)\n")
	}

	// Exits
	log.Printf("R: %v", r.Exits)
	builder.WriteString("|- Exits:\n")
	if len(r.Exits) == 0 {
		builder.WriteString("|-- (No obvious exits)\n")
	} else {
		for dir := range r.Exits {
			builder.WriteString("|-- " + dir + "\n")
		}
	}

	return builder.String()
}
