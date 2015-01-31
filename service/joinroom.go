package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/common/serviceerror"
	"github.com/protogalaxy/service-goroom/lobby"
	"golang.org/x/net/context"
)

type JoinRoomRequest struct {
	UserID string `json:"user_id"`
	RoomID string `json:"-"`
}

type JoinRoomResponse struct{}

type JoinRoom struct {
	Lobby lobby.Lobby
}

func decodeRequest(body io.Reader, r *JoinRoomRequest) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&r)
	if err != nil {
		return serviceerror.BadRequest("Unable to decode request body", err)
	}

	if r.UserID == "" {
		return serviceerror.BadRequest("Missing or empty user id", nil)
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

	// TODO: handle error
	h.Lobby.JoinRoom(req.RoomID, req.UserID)

	res := JoinRoomResponse{}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&res)
	if err != nil {
		return serviceerror.InternalServerError("Problem encoding response", err)
	}
	return nil
}

// Do implements saola.Service.
func (h *JoinRoom) Do(ctx context.Context) error {
	return httpservice.Do(h, ctx)
}

func (h *JoinRoom) Name() string {
	return "createroom"
}
