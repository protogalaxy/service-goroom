package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/arjantop/saola/httpservice"
	"github.com/golang/glog"
	"github.com/protogalaxy/common/serviceerror"
	"github.com/protogalaxy/service-goroom/lobby"
	"golang.org/x/net/context"
)

type JoinRoomRequest struct {
	UserID string `json:"user_id"`
	RoomID string `json:"-"`
}

type JoinRoomResponse struct {
	Status string `json:"status"`
}

type JoinRoom struct {
	Lobby lobby.Lobby
}

func decodeRequest(body io.Reader, r *JoinRoomRequest) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&r)
	if err != nil {
		serr := serviceerror.BadRequest("invalid_request", "Unable to decode request body")
		serr.Cause = err
		return serr
	}

	if r.UserID == "" {
		return serviceerror.BadRequest("invalid_request", "Missing or empty user id")
	}
	return nil
}

// DoHTTP implements saola.HttpService.
func (h *JoinRoom) DoHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req JoinRoomRequest
	req.RoomID = httpservice.GetParams(ctx).Get("roomID")
	err := decodeRequest(r.Body, &req)
	if err != nil {
		return err
	}

	err = h.Lobby.JoinRoom(req.RoomID, req.UserID)
	res := JoinRoomResponse{
		Status: "joined",
	}

	if err == lobby.ErrRoomNotFound {
		res.Status = "room_not_found"
	} else if err == lobby.ErrAlreadyInRoom {
		res.Status = "already_in_room"
	} else if err == lobby.ErrRoomFull {
		res.Status = "room_full"
	} else if err != nil {
		glog.Errorf("Unexpected error: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&res)
	if err != nil {
		return serviceerror.InternalServerError("server_error", "Problem encoding response", err)
	}
	return nil
}

// Do implements saola.Service.
func (h *JoinRoom) Do(ctx context.Context) error {
	return httpservice.Do(h, ctx)
}

func (h *JoinRoom) Name() string {
	return "joinroom"
}
