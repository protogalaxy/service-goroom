package lobby

import (
	"errors"
	"sync"

	"code.google.com/p/go-uuid/uuid"
)

type Room struct {
	ID          uuid.UUID
	Owner       string
	OtherPlayer string
}

var ErrRoomFull = errors.New("room full")

func (r *Room) Join(userID string) error {
	if r.OtherPlayer != "" {
		return ErrRoomFull
	}
	r.OtherPlayer = userID
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
	ErrRoomNotFound  = errors.New("room not found")
	ErrAlreadyInRoom = errors.New("already in room")
)

func (l *Manager) JoinRoom(roomID, userID string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	room, ok := l.rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}

	for _, room := range l.rooms {
		if userInRoom(room, userID) {
			return ErrAlreadyInRoom
		}
	}

	if err := room.Join(userID); err != nil {
		return err
	}
	return nil
}

func userInRoom(r *Room, userID string) bool {
	return r.Owner == userID || r.OtherPlayer == userID
}
