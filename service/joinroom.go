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
	"net/http"

	"github.com/arjantop/saola/httpservice"
	"golang.org/x/net/context"
)

type JoinRoom struct {
	Service *GoRoom
}

// DoHTTP implements saola.HttpService.
func (h *JoinRoom) DoHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req JoinRoomRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		return err
	}

	req.RoomID = httpservice.GetParams(ctx).Get("roomID")

	var res JoinRoomResponse
	err = h.Service.JoinRoom(&req, &res)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &res)
}

// Do implements saola.Service.
func (h *JoinRoom) Do(ctx context.Context) error {
	return httpservice.Do(h, ctx)
}

// Name implements saola.Service.
func (h *JoinRoom) Name() string {
	return "joinroom"
}
