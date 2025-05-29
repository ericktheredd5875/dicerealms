package game

import (
	"testing"
)

func TestRoll4d6DropLowest(t *testing.T) {
	for i := 0; i < 100; i++ {
		total, rolls := roll4d6DropLowest()
		if len(rolls) != 4 {
			t.Errorf("Expected 4 rolls, got %d", len(rolls))
		}
		for _, roll := range rolls {
			if roll < 1 || roll > 6 {
				t.Errorf("Invalid die roll: %d", roll)
			}
		}
		if total < 3 || total > 18 {
			t.Errorf("Total out of bounds: got %d", total)
		}
	}
}

func TestAssignStat(t *testing.T) {
	p := &Player{AssignedStats: make(map[string]bool)}
	stats := []string{"STR", "DEX", "CON", "INT", "WIS", "CHA"}

	for _, stat := range stats {
		msg, err := p.AssignStat(stat)
		if err != nil {
			t.Errorf("Unexpected error assigning %s: %v", stat, err)
		}
		if msg == "" {
			t.Errorf("Empty message for %s assignment", stat)
		}
		if !p.AssignedStats[stat] {
			t.Errorf("%s should be marked as assigned", stat)
		}

		// Re-assign should error
		_, err = p.AssignStat(stat)
		if err == nil {
			t.Errorf("Expected error when re-assigning %s, got none", stat)
		}
	}
}

func TestAutoGenerateStats(t *testing.T) {
	p := &Player{AssignedStats: make(map[string]bool)}
	result, err := p.AutoGenStats()
	if err != nil {
		t.Errorf("Unexpected error in AutoGenerateStats: %v", err)
	}
	if result == "" {
		t.Error("Expected non-empty result from AutoGenerateStats")
	}

	// Should fail on second call
	_, err = p.AutoGenStats()
	if err == nil {
		t.Error("Expected error when calling AutoGenerateStats again")
	}
}
