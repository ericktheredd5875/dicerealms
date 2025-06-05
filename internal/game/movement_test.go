package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerMoveSuccess(t *testing.T) {
	start := &Room{
		Name:    "Start",
		Exits:   make(map[string]*Room),
		Players: make(map[string]*Player),
	}
	target := &Room{
		Name:    "Forest",
		Exits:   make(map[string]*Room),
		Players: make(map[string]*Player),
	}
	start.Exits["north"] = target

	player := &Player{Name: "Hero", Room: start}
	start.AddPlayer(player)

	msg, err := player.Move("north")
	assert.NoError(t, err)
	assert.Contains(t, msg, "You move north into Forest")
	assert.Equal(t, target, player.Room)
	_, inStart := start.Players["Hero"]
	_, inTarget := target.Players["Hero"]
	assert.False(t, inStart)
	assert.True(t, inTarget)
}

func TestPlayerMoveInvalidDirection(t *testing.T) {
	start := &Room{
		Name:    "Start",
		Exits:   make(map[string]*Room),
		Players: make(map[string]*Player),
	}
	player := &Player{Name: "Hero", Room: start}
	start.AddPlayer(player)

	msg, err := player.Move("west")
	assert.Error(t, err)
	assert.Equal(t, "", msg)
	assert.Equal(t, start, player.Room)
	_, stillThere := start.Players["Hero"]
	assert.True(t, stillThere)
}
