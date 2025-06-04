package db

import (
	"errors"
	"log"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

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
