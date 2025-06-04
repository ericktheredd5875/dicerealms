package db

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PlayerModel struct {
	gorm.Model
	PublicID   string         `gorm:"uniqueIndex:not null"`
	Name       string         `gorm:"not null"`
	STR        int            `gorm:"default:0"`
	DEX        int            `gorm:"default:0"`
	CON        int            `gorm:"default:0"`
	INT        int            `gorm:"default:0"`
	WIS        int            `gorm:"default:0"`
	CHA        int            `gorm:"default:0"`
	Inventory  pq.StringArray `gorm:"type:text[]"`
	Equipped   pq.StringArray `gorm:"type:text[]"`
	Gold       int            `gorm:"default:0"`
	Level      int            `gorm:"default:1"`
	XP         int            `gorm:"default:0"`
	RoomID     int
	LastRoomID int
}

type RoomModel struct {
	gorm.Model
	RoomName    string `gorm:"not null"`
	Description string
	Exits       pq.StringArray `gorm:"type:text[]"`
	JoinMsg     string
	LeaveMsg    string
}

type ItemModel struct {
	gorm.Model
	Name        string
	Description string
	RoomID      int
	RoomFound   RoomModel `gorm:"foreignKey:RoomID"`
}

type SceneModel struct {
	gorm.Model
	Title       string
	Description string
	Mood        string
	RoomID      int
	Room        RoomModel `gorm:"foreignKey:RoomID"`
	StartedByID int
	StartedBy   PlayerModel `gorm:"foreignKey:StartedByID"`
}

type SceneLogModel struct {
	gorm.Model
	SceneLogID int
	Scene      SceneModel `gorm:"foreignKey:SceneLogID"`
	PlayerID   int
	Player     PlayerModel `gorm:"foreignKey:PlayerID"`
	Action     string
}
