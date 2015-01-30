package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/protogalaxy/service-goroom/service"
	"golang.org/x/net/context"
)

type LobbyMock struct {
	OnCreateRoom func(userID string) string
}

func (m LobbyMock) CreateRoom(userID string) string {
	return m.OnCreateRoom(userID)
}

func TestCreateRoomSuccess(t *testing.T) {
	l := LobbyMock{
		OnCreateRoom: func(userID string) string {
			return "room1"
		},
	}
	s := service.CreateRoom{
		Lobby: l,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader(`{"user_id": "user123"}`))
	err := s.DoHTTP(context.Background(), w, req)
	if err != nil {
		t.Fatalf("Creating a room should no fail but got: %s", err)
	}
	if w.Code != http.StatusCreated {
		t.Errorf("Should respond with status 'Created' but got: %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Unexpected content type: %s", w.Header().Get("Content-Type"))
	}
	res := make(map[string]interface{})
	expected := map[string]interface{}{
		"room_id": "room1",
	}
	dec := json.NewDecoder(w.Body)
	if err := dec.Decode(&res); err != nil || !reflect.DeepEqual(res, expected) {
		t.Errorf("Invalid response body: expected '%v' but got '%v'", expected, res)
		if err != nil {
			t.Logf("Error: %s", err)
		}
	}
}

func TestCreateRoomDecodeBodyError(t *testing.T) {
	s := service.CreateRoom{
		Lobby: LobbyMock{},
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader(`{`))
	err := s.DoHTTP(context.Background(), w, req)
	if err == nil {
		t.Fatal("Expecting error but got none")
	}
}

func TestCreateRoomMissingRoomId(t *testing.T) {
	s := service.CreateRoom{
		Lobby: LobbyMock{},
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader(`{}`))
	err := s.DoHTTP(context.Background(), w, req)
	if err == nil {
		t.Fatal("Expecting error but got none")
	}
}
