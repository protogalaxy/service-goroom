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

type GoRoom struct {
	Lobby lobby.Lobby
}

type CreateRoomRequest struct {
	UserID string `json:"user_id"`
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

func (s *GoRoom) CreateRoom(req *CreateRoomRequest, res *CreateRoomResponse) error {
	if req.UserID == "" {
		return NewMissingParameterError("user id")
	}

	res.RoomID = s.Lobby.CreateRoom(req.UserID)
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
		res.Status = "room_not_found"
	} else if err == lobby.ErrAlreadyInRoom {
		res.Status = "already_in_room"
	} else if err == lobby.ErrRoomFull {
		res.Status = "room_full"
	} else if err != nil {
		glog.Errorf("Unexpected error: %s", err)
		return err
	} else {
		res.Status = "joined"
	}
	return nil
}
