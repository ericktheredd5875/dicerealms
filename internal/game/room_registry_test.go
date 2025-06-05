package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetRoomByName(t *testing.T) {
	roomA := &Room{Name: "Tavern"}
	roomB := &Room{Name: "Dungeon"}

	SetRooms(map[string]*Room{
		"Tavern":  roomA,
		"Dungeon": roomB,
	})

	result := GetRoomByName("Tavern")
	assert.NotNil(t, result)
	assert.Equal(t, "Tavern", result.Name)

	nonexistent := GetRoomByName("Castle")
	assert.Nil(t, nonexistent)
}

func TestGetAllRooms(t *testing.T) {
	roomX := &Room{Name: "Garden"}
	roomY := &Room{Name: "Tower"}

	SetRooms(map[string]*Room{
		"Garden": roomX,
		"Tower":  roomY,
	})

	allRooms := GetAllRooms()
	assert.Len(t, allRooms, 2)
	assert.Contains(t, allRooms, "Garden")
	assert.Contains(t, allRooms, "Tower")
	assert.Equal(t, roomX, allRooms["Garden"])
	assert.Equal(t, roomY, allRooms["Tower"])
}
