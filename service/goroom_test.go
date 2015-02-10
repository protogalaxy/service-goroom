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

package service_test

import (
	"strings"
	"testing"

	"github.com/protogalaxy/service-goroom/lobby"
	"github.com/protogalaxy/service-goroom/service"
)

type LobbyMock struct {
	OnCreateRoom func(userID string) string
	OnJoinRoom   func(roomID, userID string) error
}

func (m LobbyMock) CreateRoom(userID string) string {
	return m.OnCreateRoom(userID)
}
func (m LobbyMock) JoinRoom(roomID, userID string) error {
	return m.OnJoinRoom(roomID, userID)
}

func TestGoRoomIsCreated(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnCreateRoom: func(userID string) string {
				if userID != "userid" {
					t.Errorf("Unexpected user id: %s", userID)
				}
				return "roomid"
			},
		},
	}
	req := service.CreateRoomRequest{
		UserID: "userid",
	}
	var res service.CreateRoomResponse
	err := s.CreateRoom(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if res.RoomID != "roomid" {
		t.Errorf("Invalid room id: %s", res.RoomID)
	}
}

func AssertMissingParameterError(t *testing.T, err error, expected string) {
	mpe, ok := err.(service.MissingParameterError)
	if !ok {
		t.Fatalf("Invalid error type: %#v", err)
	}
	if !strings.Contains(mpe.Error(), expected) {
		t.Errorf("Invalid error parameter error: %s", mpe)
	}
}

func TestGoRoomCreateErrorReturnedIfUserIDMissing(t *testing.T) {
	s := &service.GoRoom{}
	req := service.CreateRoomRequest{}
	var res service.CreateRoomResponse
	AssertMissingParameterError(t, s.CreateRoom(&req, &res), "user id")
}

func TestGoRoomUserCanJoinTheRoom(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnJoinRoom: func(roomID, userID string) error {
				if roomID != "roomid" {
					t.Errorf("Unexpected room id: %s", userID)
				}
				if userID != "userid" {
					t.Errorf("Unexpected user id: %s", userID)
				}
				return nil
			},
		},
	}
	req := service.JoinRoomRequest{
		RoomID: "roomid",
		UserID: "userid",
	}
	var res service.JoinRoomResponse
	err := s.JoinRoom(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if res.Status != "joined" {
		t.Errorf("Unexpected room join status: %s", res.Status)
	}
}

func TestGoRoomUserCanJoin(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnJoinRoom: func(roomID, userID string) error {
				if roomID != "roomid" {
					t.Errorf("Unexpected room id: %s", userID)
				}
				if userID != "userid" {
					t.Errorf("Unexpected user id: %s", userID)
				}
				return nil
			},
		},
	}
	req := service.JoinRoomRequest{
		RoomID: "roomid",
		UserID: "userid",
	}
	var res service.JoinRoomResponse
	err := s.JoinRoom(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if res.Status != "joined" {
		t.Errorf("Unexpected room join status: %s", res.Status)
	}
}

func TestGoRoomJoinErrorIfRoomIDMissing(t *testing.T) {
	s := &service.GoRoom{}
	req := service.JoinRoomRequest{
		UserID: "userid",
	}
	var res service.JoinRoomResponse
	AssertMissingParameterError(t, s.JoinRoom(&req, &res), "room id")
}

func TestGoRoomJoinErrorIfUserIDMissing(t *testing.T) {
	s := &service.GoRoom{}
	req := service.JoinRoomRequest{
		RoomID: "roomid",
	}
	var res service.JoinRoomResponse
	AssertMissingParameterError(t, s.JoinRoom(&req, &res), "user id")
}

func TestGoRoomUserCannotJoinTheNonexistentRoom(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnJoinRoom: func(roomID, userID string) error {
				return lobby.ErrRoomNotFound
			},
		},
	}
	req := service.JoinRoomRequest{
		RoomID: "nonexistent",
		UserID: "userid",
	}
	var res service.JoinRoomResponse
	err := s.JoinRoom(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if res.Status != "room_not_found" {
		t.Errorf("Unexpected room join status: %s", res.Status)
	}
}

func TestGoRoomUserCannotJoinTwoRoomsAtOnce(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnJoinRoom: func(roomID, userID string) error {
				return lobby.ErrAlreadyInRoom
			},
		},
	}
	req := service.JoinRoomRequest{
		RoomID: "roomid",
		UserID: "alreadyjoineduser",
	}
	var res service.JoinRoomResponse
	err := s.JoinRoom(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if res.Status != "already_in_room" {
		t.Errorf("Unexpected room join status: %s", res.Status)
	}
}

func TestGoRoomUserCannotJoinTheFullRoom(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnJoinRoom: func(roomID, userID string) error {
				return lobby.ErrRoomFull
			},
		},
	}
	req := service.JoinRoomRequest{
		RoomID: "fullroom",
		UserID: "userid",
	}
	var res service.JoinRoomResponse
	err := s.JoinRoom(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if res.Status != "room_full" {
		t.Errorf("Unexpected room join status: %s", res.Status)
	}
}
