package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ericktheredd5875/dicerealms/pkg/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SeedItem struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Rarity      string `json:"Rarity"`
	Effect      string `json:"Effect"`
	Category    string `json:"Category"`
	RoomFound   string `json:"RoomFound"`
}

var initRooms = []RoomModel{
	{
		RoomName:    "The Tavern",
		Description: "A cozy tavern with a roaring fire and the smell of ale.",
		Exits:       pq.StringArray{"north:Dark Forest", "east:Crystal Cave"},
		JoinMsg:     "enters the tavern, greeted by warmth and laughter.",
		LeaveMsg:    "steps out through the creaky tavern door.",
	},
	{
		RoomName:    "Dark Forest",
		Description: "Tall, ominous trees block out most of the sunlight.",
		Exits:       pq.StringArray{"south:The Tavern"},
		JoinMsg:     "emerges into the shadowy forest.",
		LeaveMsg:    "disappears into the thick, dark underbrush.",
	},
	{
		RoomName:    "Crystal Cave",
		Description: "A cavern filled with glowing crystals that hum softly.",
		Exits:       pq.StringArray{"west:The Tavern"},
		JoinMsg:     "descends into the shimmering Crystal Cave.",
		LeaveMsg:    "climbs out of the cave’s glowing passage.",
	},
}

func SeedRooms(db *gorm.DB) {
	for _, room := range initRooms {
		var existing RoomModel
		err := db.Where("room_name = ?", room.RoomName).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&room).Error; err != nil {
				log.Printf("Failed to seed room %s: %v", room.RoomName, err)
			}
		}
	}
}

func LoadItems(db *gorm.DB, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read items file: %w", err)
	}

	var items []SeedItem
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to unmarshal items: %w", err)
	}

	for _, i := range items {

		ItemID, _ := utils.GenerateKHash(i.Name+":DiceRealms:Items:4000", "")
		item := ItemModel{
			ItemID:      ItemID,
			Name:        i.Name,
			Description: i.Description,
			Rarity:      i.Rarity,
			Effect:      i.Effect,
			Category:    i.Category,
			RoomFoundID: 0,
		}

		if i.RoomFound != "" {
			var room RoomModel
			if err := db.Where("room_name = ?", i.RoomFound).First(&room).Error; err != nil {
				item.RoomFoundID = room.ID
			}
		}

		var existing ItemModel
		if err := db.Where("name = ?", item.Name).First(&existing).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&item)
		}
	}

	log.Printf("Loaded %d items", len(items))
	return nil
}
