// Copyright (C) 2015 The Protogalaxy Project
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lobby

import (
	"errors"
	"sync"

	"code.google.com/p/go-uuid/uuid"
)

type Room struct {
	ID          string
	Owner       string
	OtherPlayer string
}

var (
	ErrRoomFull      = errors.New("room full")
	ErrRoomNotFound  = errors.New("room not found")
	ErrAlreadyInRoom = errors.New("already in room")
)

func (r *Room) Join(userID string) error {
	if r.OtherPlayer != "" {
		return ErrRoomFull
	}
	r.OtherPlayer = userID
	return nil
}

type Lobby interface {
	RoomInfo(roomID string) (*Room, error)
	CreateRoom(userID string) (string, error)
	JoinRoom(roomID, userID string) error
}

type Generator interface {
	GenerateID() string
}

type UUIDGenerator struct{}

func (g *UUIDGenerator) GenerateID() string {
	return uuid.NewRandom().String()
}

type Manager struct {
	lock  sync.Mutex
	rooms map[string]*Room

	Generator Generator
}

func NewLobby() *Manager {
	return &Manager{
		rooms:     make(map[string]*Room),
		Generator: &UUIDGenerator{},
	}
}

func (l *Manager) CreateRoom(userID string) (string, error) {
	if l.isUserInAnyRoom(userID) {
		return "", ErrAlreadyInRoom
	}

	r := &Room{
		ID:    l.Generator.GenerateID(),
		Owner: userID,
	}
	l.lock.Lock()
	l.rooms[r.ID] = r
	l.lock.Unlock()
	return r.ID, nil
}

func (l *Manager) RoomInfo(roomID string) (*Room, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	room, ok := l.rooms[roomID]
	if !ok {
		return nil, ErrRoomNotFound
	}

	return &Room{
		ID:          room.ID,
		Owner:       room.Owner,
		OtherPlayer: room.OtherPlayer,
	}, nil
}

func (l *Manager) JoinRoom(roomID, userID string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	room, ok := l.rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}

	if l.isUserInAnyRoom(userID) {
		return ErrAlreadyInRoom
	}

	if err := room.Join(userID); err != nil {
		return err
	}
	return nil
}

func (l *Manager) isUserInAnyRoom(userID string) bool {
	for _, room := range l.rooms {
		if userInRoom(room, userID) {
			return true
		}
	}
	return false
}

func userInRoom(r *Room, userID string) bool {
	return r.Owner == userID || r.OtherPlayer == userID
}
