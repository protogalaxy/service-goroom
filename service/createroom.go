package service

import (
	"encoding/json"
	"net/http"

	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/common/serviceerror"
	"github.com/protogalaxy/service-goroom/lobby"
	"golang.org/x/net/context"
)

type CreateRoomRequest struct {
	UserID string `json:"user_id"`
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

type CreateRoom struct {
	Lobby lobby.Lobby
}

// DoHTTP implements saola.HttpService.
func (h *CreateRoom) DoHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var cr CreateRoomRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cr)
	if err != nil {
		return serviceerror.BadRequest("Unable to decode request body", err)
	}

	if cr.UserID == "" {
		return serviceerror.BadRequest("Missing user id", nil)
	}

	res := CreateRoomResponse{
		RoomID: h.Lobby.CreateRoom(cr.UserID),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&res)
	if err != nil {
		return serviceerror.InternalServerError("Problem encoding response", err)
	}
	return nil
}

// Do implements saola.Service.
func (h *CreateRoom) Do(ctx context.Context) error {
	return httpservice.Do(h, ctx)
}

func (h *CreateRoom) Name() string {
	return "createroom"
}
