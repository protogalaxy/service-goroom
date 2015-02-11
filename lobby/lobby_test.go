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

package lobby_test

import (
	"testing"

	"github.com/protogalaxy/service-goroom/lobby"
)

type MockGenerator struct {
	OnGenerateID func() string
}

func (m *MockGenerator) GenerateID() string {
	return m.OnGenerateID()
}

func TestGeneratorGeneratedUniqueValues(t *testing.T) {
	g := lobby.UUIDGenerator{}
	v1 := g.GenerateID()
	if v1 == "" {
		t.Errorf("Generated id should not be empty")
	}
	v2 := g.GenerateID()
	if v2 == "" {
		t.Errorf("Generated id should not be empty")
	}
	if g.GenerateID() == g.GenerateID() {
		t.Errorf("Two generated ids should not be the same: %s == %s", v1, v2)
	}
}

func TestLobbyIDOfCreatedRoomIsReturned(t *testing.T) {
	l := lobby.NewLobby()
	l.Generator = &MockGenerator{
		OnGenerateID: func() string {
			return "roomid"
		},
	}
	rid, _ := l.CreateRoom("userid")
	if rid != "roomid" {
		t.Errorf("Unexpected room id: %s", rid)
	}
}

func TestLobbyUserCannotCreateNewRoomIfAlreadyJoinedAnotherRoom(t *testing.T) {
	l := lobby.NewLobby()

	rid, _ := l.CreateRoom("user1")
	l.JoinRoom(rid, "user2")

	_, err := l.CreateRoom("user2")
	if err != lobby.ErrAlreadyInRoom {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

func TestLobbyUserCannotCreateNewRoomIfAlreadyOwnerOfAnotherRoom(t *testing.T) {
	l := lobby.NewLobby()

	l.CreateRoom("user1")

	_, err := l.CreateRoom("user1")
	if err != lobby.ErrAlreadyInRoom {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

func TestLobbyCanRetrieveCreatedRoomInfo(t *testing.T) {
	l := lobby.NewLobby()
	l.Generator = &MockGenerator{
		OnGenerateID: func() string {
			return "roomid"
		},
	}

	rid, _ := l.CreateRoom("userid")

	rinfo, err := l.RoomInfo(rid)
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
	expected := lobby.Room{
		ID:    "roomid",
		Owner: "userid",
	}
	if *rinfo != expected {
		t.Errorf("Invalid response: %#v != %#v", rinfo, expected)
	}
}

func TestLobbyCannotReceiveInfoOfNonexistentRoom(t *testing.T) {
	l := lobby.NewLobby()

	_, err := l.RoomInfo("roomid")
	if err != lobby.ErrRoomNotFound {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

func TestLobbyCreatedRoomCanBeJoined(t *testing.T) {
	l := lobby.NewLobby()
	l.Generator = &MockGenerator{
		OnGenerateID: func() string {
			return "roomid"
		},
	}

	rid, _ := l.CreateRoom("userid")
	err := l.JoinRoom(rid, "user2")
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	rinfo, err := l.RoomInfo(rid)
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
	expected := lobby.Room{
		ID:          "roomid",
		Owner:       "userid",
		OtherPlayer: "user2",
	}
	if *rinfo != expected {
		t.Errorf("Invalid response: %#v != %#v", rinfo, expected)
	}
}

func TestLobbyUserCannotJoinNonexistentRoom(t *testing.T) {
	l := lobby.NewLobby()

	err := l.JoinRoom("roomid", "user2")
	if err != lobby.ErrRoomNotFound {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

func TestLobbyUserCanOnlyJoinOneRoomAtOnce(t *testing.T) {
	l := lobby.NewLobby()

	rid, _ := l.CreateRoom("user1")
	l.JoinRoom(rid, "user2")

	rid2, _ := l.CreateRoom("user3")

	// User is the other player int he room
	err := l.JoinRoom(rid2, "user2")
	if err != lobby.ErrAlreadyInRoom {
		t.Fatalf("Unexpected error: %#v", err)
	}

	// User is the owner of the room
	err = l.JoinRoom(rid2, "user1")
	if err != lobby.ErrAlreadyInRoom {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

func TestLobbyUserCannotJoinFullRoom(t *testing.T) {
	l := lobby.NewLobby()

	rid, _ := l.CreateRoom("userid")
	l.JoinRoom(rid, "user2")

	err := l.JoinRoom(rid, "user3")
	if err != lobby.ErrRoomFull {
		t.Fatalf("Unexpected error: %#v", err)
	}
}
