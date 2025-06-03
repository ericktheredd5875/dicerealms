package session

import "sync"

var (
	activeSessions = make(map[*Session]bool)
	sessionMu      sync.Mutex
)

func RegisterSession(s *Session) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	activeSessions[s] = true
}

func UnregisterSession(s *Session) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	delete(activeSessions, s)
}

func SaveAllSessions() {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	for s := range activeSessions {
		if s.Player != nil {
			_ = s.Player.Save()
		}

		s.Close()
	}
}
