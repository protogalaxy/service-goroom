package lobby

import (
	"errors"
	"sync"

	"code.google.com/p/go-uuid/uuid"
)

type Room struct {
	ID    uuid.UUID
	Owner string
}

func (r *Room) Join(userID string) error {
	return nil
}

type Lobby interface {
	CreateRoom(userID string) string
	JoinRoom(roomID, userID string) error
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

var (
	ErrRoomDoesNotExist = errors.New("Room does not exist")
)

func (l *Manager) JoinRoom(roomID, userID string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	room, ok := l.rooms[roomID]
	if !ok {
		return ErrRoomDoesNotExist
	}
	if err := room.Join(userID); err != nil {
		return err
	}
	return nil
}
