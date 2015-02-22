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

package service

import (
	"github.com/golang/glog"
	"github.com/protogalaxy/service-goroom/lobby"
)

type MissingParameterError string

func NewMissingParameterError(name string) MissingParameterError {
	return MissingParameterError(name)
}

func (e MissingParameterError) Error() string {
	return "Missing parameter: " + string(e)
}

const (
	RoomStatusCreated       = "created"
	RoomStatusAlreadyInRoom = "already_in_room"
	RoomStatusFull          = "room_full"
	RoomStatusFound         = "found"
	RoomStatusNotFound      = "not_found"
	RoomStatusJoined        = "joined"
)

type GoRoom struct {
	Lobby lobby.Lobby
}

type CreateRoomRequest struct {
	UserID string `json:"user_id"`
}

type CreateRoomResponse struct {
	Status string `json:"status"`
	RoomID string `json:"room_id"`
}

func (s *GoRoom) CreateRoom(req *CreateRoomRequest, res *CreateRoomResponse) error {
	if req.UserID == "" {
		return NewMissingParameterError("user id")
	}

	rid, err := s.Lobby.CreateRoom(req.UserID)
	if err == lobby.ErrAlreadyInRoom {
		res.Status = RoomStatusAlreadyInRoom
		return nil
	} else if err != nil {
		glog.Errorf("Unexpected error: %s", err)
		return err
	}
	res.Status = RoomStatusCreated
	res.RoomID = rid
	return nil
}

type JoinRoomRequest struct {
	UserID string `json:"user_id"`
	RoomID string `json:"-"`
}

type JoinRoomResponse struct {
	Status string `json:"status"`
}

func (s *GoRoom) JoinRoom(req *JoinRoomRequest, res *JoinRoomResponse) error {
	if req.UserID == "" {
		return NewMissingParameterError("user id")
	} else if req.RoomID == "" {
		return NewMissingParameterError("room id")
	}

	err := s.Lobby.JoinRoom(req.RoomID, req.UserID)

	if err == lobby.ErrRoomNotFound {
		res.Status = RoomStatusNotFound
	} else if err == lobby.ErrAlreadyInRoom {
		res.Status = RoomStatusAlreadyInRoom
	} else if err == lobby.ErrRoomFull {
		res.Status = RoomStatusFull
	} else if err != nil {
		glog.Errorf("Unexpected error: %s", err)
		return err
	} else {
		res.Status = RoomStatusJoined
	}
	return nil
}

type RoomInfoRequest struct {
	RoomID string `json:"-"`
}

type Room struct {
	RoomID      string `json:"room_id"`
	Owner       string `json:"owner"`
	OtherPlayer string `json:"other_player,omitempty"`
}

type RoomInfoResponse struct {
	Status string `json:"status"`
	Room   *Room  `json:"room"`
}

func (s *GoRoom) RoomInfo(req *RoomInfoRequest, res *RoomInfoResponse) error {
	room, err := s.Lobby.RoomInfo(req.RoomID)
	if err == lobby.ErrRoomNotFound {
		res.Status = RoomStatusNotFound
	} else if err != nil {
		glog.Errorf("Unexpected error: %s", err)
		return err
	} else {
		res.Status = RoomStatusFound
		res.Room = &Room{
			RoomID:      room.ID,
			Owner:       room.Owner,
			OtherPlayer: room.OtherPlayer,
		}
	}
	return nil
}
