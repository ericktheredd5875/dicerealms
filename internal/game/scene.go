package game

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Scene struct {
	Title     string
	Mood      string
	StartedBy string
	StartedAt time.Time
	EndedBy   string
	EndedAt   time.Time
	Log       []string
}

func (r *Room) StartScene(title string, mood, startedBy string) {
	r.mu.Lock()

	r.ActiveScene = &Scene{
		Title:     title,
		Mood:      mood,
		StartedBy: startedBy,
		StartedAt: time.Now(),
		Log:       []string{},
	}
	r.mu.Unlock()

	r.Broadcast(fmt.Sprintf("<new> Scene started: %s [%s]", title, mood), "")
}

func (r *Room) EndScene(endedBy string) string {
	r.mu.Lock()

	if r.ActiveScene == nil {
		return "No active scene to end."
	}

	r.ActiveScene.EndedBy = endedBy
	r.ActiveScene.EndedAt = time.Now()
	summary := r.ActiveScene.Summary()

	filename := fmt.Sprintf("scene_%s_%s.txt", sanitizeFileName(r.ActiveScene.Title), time.Now().Format("20060102_150405"))
	path := filepath.Join("scenes", filename)
	_ = os.MkdirAll("scenes", 0755)
	err := os.WriteFile(path, []byte(summary), 0644)
	if err != nil {
		summary += fmt.Sprintf("Error writing scene file: %v", err)
	} else {
		summary += fmt.Sprintf("Scene saved to %s", path)
	}

	r.mu.Unlock()

	r.ActiveScene = nil
	r.Broadcast("<end> Scene ended", "")

	return summary

}

func (s *Scene) LogEntry(msg string) {
	s.Log = append(s.Log, fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg))
}

func (s *Scene) Summary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Scene: %s\nMood: %s\nStarted By: %s at %s\n\n", s.Title, s.Mood, s.StartedBy, s.StartedAt.Format(time.Kitchen)))

	sb.WriteString("Log:\n")
	for _, entry := range s.Log {
		sb.WriteString(fmt.Sprintf(entry + "\n"))
	}

	if s.EndedBy != "" {
		sb.WriteString(fmt.Sprintf("Ended By: %s at %s\n", s.EndedBy, s.EndedAt.Format(time.Kitchen)))
	}

	return sb.String()
}

func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(" ", "_", "/", "-", "\\", "-", ":", "_")
	return replacer.Replace(name)
}
