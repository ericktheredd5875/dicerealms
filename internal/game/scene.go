package game

import (
	"fmt"
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
