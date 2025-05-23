package game

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var diceExpr = regexp.MustCompile(`^(\d+)d(\d+)([+-]\d+)?$`)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Roll(expr string) (int, string, error) {
	matches := diceExpr.FindStringSubmatch(expr)
	if matches == nil {
		return 0, "", fmt.Errorf("invalid dice expression: %q", expr)
	}

	numDice, _ := strconv.Atoi(matches[1])
	diceSides, _ := strconv.Atoi(matches[2])

	modifier := 0
	if matches[3] != "" {
		modifier, _ = strconv.Atoi(matches[3])
	}

	rolls := make([]int, numDice)
	total := 0
	for i := 0; i < numDice; i++ {
		roll := rand.Intn(diceSides) + 1
		rolls[i] = roll
		total += roll
	}

	total += modifier
	parts := make([]string, len(rolls))
	for i, r := range rolls {
		parts[i] = strconv.Itoa(r)
	}

	details := fmt.Sprintf("Rolled %s -> [%s]+%s= %d",
		expr, strings.Join(parts, ", "),
		fmtIf(modifier != 0, fmt.Sprintf("%d", modifier), ""),
		total,
	)

	return total, details, nil
}

func fmtIf(cond bool, a string, b string) string {
	if cond {
		return a
	}

	return b
}
