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
	"reflect"
	"strings"
	"testing"

	"github.com/protogalaxy/service-goroom/lobby"
	"github.com/protogalaxy/service-goroom/service"
)

type LobbyMock struct {
	OnRoomInfo   func(roomID string) (*lobby.Room, error)
	OnCreateRoom func(userID string) (string, error)
	OnJoinRoom   func(roomID, userID string) error
}

func (m LobbyMock) RoomInfo(roomID string) (*lobby.Room, error) {
	return m.OnRoomInfo(roomID)
}

func (m LobbyMock) CreateRoom(userID string) (string, error) {
	return m.OnCreateRoom(userID)
}

func (m LobbyMock) JoinRoom(roomID, userID string) error {
	return m.OnJoinRoom(roomID, userID)
}

func TestGoRoomIsCreated(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnCreateRoom: func(userID string) (string, error) {
				if userID != "userid" {
					t.Errorf("Unexpected user id: %s", userID)
				}
				return "roomid", nil
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
	expected := service.CreateRoomResponse{
		Status: "created",
		RoomID: "roomid",
	}
	if res != expected {
		t.Errorf("Invalid response: %#v != %#v", res, expected)
	}
}

func TestGoRoomUserCreatingTheRoomAlreadyInAnotherRoom(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnCreateRoom: func(userID string) (string, error) {
				return "", lobby.ErrAlreadyInRoom
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
	expected := service.CreateRoomResponse{
		Status: "already_in_room",
	}
	if res != expected {
		t.Errorf("Invalid response: %#v != %#v", res, expected)
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

func TestGoRoomIsJoined(t *testing.T) {
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
	expected := service.JoinRoomResponse{
		Status: "joined",
	}
	if res != expected {
		t.Errorf("Invalid response: %#v != %#v", res, expected)
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
	expected := service.JoinRoomResponse{
		Status: "not_found",
	}
	if res != expected {
		t.Errorf("Invalid response: %#v != %#v", res, expected)
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
	expected := service.JoinRoomResponse{
		Status: "already_in_room",
	}
	if res != expected {
		t.Errorf("Invalid response: %#v != %#v", res, expected)
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
	expected := service.JoinRoomResponse{
		Status: "room_full",
	}
	if res != expected {
		t.Errorf("Invalid response: %#v != %#v", res, expected)
	}
}

func TestGoRoomRespondsWithRoomInfoForTheRoomID(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnRoomInfo: func(roomID string) (*lobby.Room, error) {
				return &lobby.Room{
					ID:          "roomid",
					Owner:       "owner",
					OtherPlayer: "other player",
				}, nil
			},
		},
	}
	req := service.RoomInfoRequest{
		RoomID: "roomid",
	}
	var res service.RoomInfoResponse
	err := s.RoomInfo(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	expected := service.RoomInfoResponse{
		Status: "found",
		Room: &service.Room{
			RoomID:      "roomid",
			Owner:       "owner",
			OtherPlayer: "other player",
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("Invalid response expected: %#v but got: %#v", expected, res)
	}
}

func TestGoRoomRespondsWithRoomNotFoundIfThereIsNoRoomWithTheWantedID(t *testing.T) {
	s := &service.GoRoom{
		Lobby: &LobbyMock{
			OnRoomInfo: func(roomID string) (*lobby.Room, error) {
				return nil, lobby.ErrRoomNotFound
			},
		},
	}
	req := service.RoomInfoRequest{
		RoomID: "nonexistent room id",
	}
	var res service.RoomInfoResponse
	err := s.RoomInfo(&req, &res)
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	expected := service.RoomInfoResponse{
		Status: "not_found",
	}
	if res != expected {
		t.Errorf("Invalid response expected: %#v but got: %#v", expected, res)
	}
}
