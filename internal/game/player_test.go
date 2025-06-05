package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssignStatStr(t *testing.T) {
	player := &Player{AssignedStats: make(map[string]bool)}
	msg, err := player.AssignStat("STR")
	assert.NoError(t, err)
	assert.Contains(t, msg, "Rolled 4d6")
	assert.True(t, player.AssignedStats["STR"])
}

func TestAssignStatAlreadyAssigned(t *testing.T) {
	player := &Player{AssignedStats: map[string]bool{"STR": true}}
	msg, err := player.AssignStat("STR")
	assert.Error(t, err)
	assert.Empty(t, msg)
}

func TestAssignStatInvalid(t *testing.T) {
	player := &Player{AssignedStats: make(map[string]bool)}
	msg, err := player.AssignStat("XYZ")
	assert.Error(t, err)
	assert.Empty(t, msg)
}

func TestAutoGenStats(t *testing.T) {
	player := &Player{AssignedStats: make(map[string]bool)}
	msg, err := player.AutoGenStats()
	assert.NoError(t, err)
	assert.Contains(t, msg, "Rolled 4d6")
	assert.Len(t, player.AssignedStats, len(validStats))
}

func TestAutoGenStatsAlreadyAssigned(t *testing.T) {
	player := &Player{AssignedStats: map[string]bool{"STR": true}}
	msg, err := player.AutoGenStats()
	assert.Error(t, err)
	assert.Empty(t, msg)
}

func TestAddRemoveItem(t *testing.T) {
	player := &Player{Inventory: []string{}}
	player.AddItem("Potion")
	assert.Contains(t, player.Inventory, "Potion")

	removed := player.RemoveItem("Potion")
	assert.True(t, removed)
	assert.NotContains(t, player.Inventory, "Potion")

	removedAgain := player.RemoveItem("Potion")
	assert.False(t, removedAgain)
}

func TestInventoryList(t *testing.T) {
	player := &Player{Inventory: []string{}}
	assert.Equal(t, "Your inventory is empty.", player.InventoryList())

	player.Inventory = []string{"Sword", "Shield"}
	list := player.InventoryList()
	assert.Contains(t, list, "Sword")
	assert.Contains(t, list, "Shield")
}

func TestShowStats(t *testing.T) {
	player := &Player{Name: "TestHero", Stats: Stats{STR: 10, DEX: 11, CON: 12, INT: 13, WIS: 14, CHA: 15}}
	statsStr := player.ShowStats()
	assert.Contains(t, statsStr, "Stats for TestHero")
	assert.Contains(t, statsStr, "STR: 10")
	assert.Contains(t, statsStr, "CHA: 15")
}
