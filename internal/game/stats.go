package game

import (
	"sort"
)

var validStats = []string{"STR", "DEX", "CON", "INT", "WIS", "CHA"}

func roll4d6DropLowest() (int, []int) {
	rolls := []int{
		rng.Intn(6) + 1,
		rng.Intn(6) + 1,
		rng.Intn(6) + 1,
		rng.Intn(6) + 1,
	}

	sort.Ints(rolls)
	return rolls[1] + rolls[2] + rolls[3], rolls
}
