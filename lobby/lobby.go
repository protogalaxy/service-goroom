package lobby

import (
	"sync"

	"code.google.com/p/go-uuid/uuid"
)

type Room struct {
	ID    uuid.UUID
	Owner string
}

type Lobby interface {
	CreateRoom(userID string) string
}

type Manager struct {
	lock  sync.Mutex
	rooms map[string]*Room
}

func NewLobby() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

func (l *Manager) CreateRoom(userID string) string {
	id := uuid.NewRandom()
	r := &Room{
		ID:    id,
		Owner: userID,
	}
	l.lock.Lock()
	l.rooms[id.String()] = r
	l.lock.Unlock()
	return r.ID.String()
}
