package service

import (
	"github.com/golang/glog"
	"github.com/protogalaxy/service-goroom/lobby"
)

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
