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

package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/arjantop/saola"
	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/common/serviceerror"
	"github.com/protogalaxy/service-goroom/lobby"
	"github.com/protogalaxy/service-goroom/service"
)

func Filters() []saola.Filter {
	return []saola.Filter{
		httpservice.NewCancellationFilter(),
		serviceerror.NewErrorResponseFilter(),
		serviceerror.NewErrorLoggerFilter(),
	}
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	lobby := lobby.NewLobby()
	goroom := &service.GoRoom{
		Lobby: lobby,
	}

	endpoint := httpservice.NewEndpoint()
	endpoint.POST("/rooms", saola.Apply(
		&service.CreateRoom{
			Service: goroom,
		},
		Filters()...))

	endpoint.POST("/rooms/:roomID/join", saola.Apply(
		&service.JoinRoom{
			Service: goroom,
		},
		Filters()...))

	endpoint.GET("/rooms/:roomID", saola.Apply(
		&service.RoomInfo{
			Service: goroom,
		},
		Filters()...))

	log.Fatal(httpservice.Serve(":10200", saola.Apply(
		endpoint,
		httpservice.NewStdRequestLogFilter())))
}
