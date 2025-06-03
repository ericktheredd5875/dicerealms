package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PlayerModel struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex"`
	STR       int
	DEX       int
	CON       int
	INT       int
	WIS       int
	CHA       int
	Inventory pq.StringArray `gorm:"type:text[]"`
	Equipped  pq.StringArray `gorm:"type:text[]"`
	Gold      int
	Level     int
	XP        int
	RoomName  string
}
