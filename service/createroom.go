package service

import (
	"encoding/json"
	"net/http"

	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/common/serviceerror"
	"golang.org/x/net/context"
)

type CreateRoom struct {
	Service *GoRoom
}

// DoHTTP implements saola.HttpService.
func (h *CreateRoom) DoHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req CreateRoomRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		serr := serviceerror.BadRequest("invalid_request", "Unable to decode request body")
		serr.Cause = err
		return serr
	}

	if req.UserID == "" {
		return serviceerror.BadRequest("invalid_request", "Missing user id")
	}

	var res CreateRoomResponse
	err = h.Service.CreateRoom(&req, &res)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&res)
	if err != nil {
		return serviceerror.InternalServerError("server_error", "Problem encoding response", err)
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
