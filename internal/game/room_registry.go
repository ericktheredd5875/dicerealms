package game

import (
	"log"
	"sync"
)

var (
	roomRegistry   = make(map[string]*Room)
	roomRegistryMu sync.RWMutex
)

func SetRooms(rooms map[string]*Room) {
	roomRegistryMu.Lock()
	defer roomRegistryMu.Unlock()
	roomRegistry = rooms
}

func GetRoomByName(name string) *Room {
	roomRegistryMu.RLock()
	defer roomRegistryMu.RUnlock()
	room, ok := roomRegistry[name]
	if !ok {
		log.Printf("Room %s not found", name)
		return nil
	}
	return room
}

func GetAllRooms() map[string]*Room {
	roomRegistryMu.RLock()
	defer roomRegistryMu.RUnlock()

	copy := make(map[string]*Room)
	for k, v := range roomRegistry {
		copy[k] = v
	}

	return copy
}
