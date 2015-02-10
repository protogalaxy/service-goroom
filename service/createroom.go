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
	"encoding/json"
	"net/http"

	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/common/serviceerror"
	"golang.org/x/net/context"
)

type CreateRoom struct {
	Service *GoRoom
}

func DecodeRequest(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(req)
	if err != nil {
		serr := serviceerror.BadRequest("invalid_request", "Unable to decode request body")
		serr.Cause = err
		return serr
	}
	return nil
}

func EncodeResponse(w http.ResponseWriter, status int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		return serviceerror.InternalServerError("server_error", "Problem encoding response", err)
	}
	return nil
}

// DoHTTP implements saola.HttpService.
func (h *CreateRoom) DoHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req CreateRoomRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		return err
	}

	var res CreateRoomResponse
	err = h.Service.CreateRoom(&req, &res)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusCreated, &res)
}

// Do implements saola.Service.
func (h *CreateRoom) Do(ctx context.Context) error {
	return httpservice.Do(h, ctx)
}

func (h *CreateRoom) Name() string {
	return "createroom"
}
