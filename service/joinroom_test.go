package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/arjantop/saola/httpservice"
	"github.com/protogalaxy/service-goroom/service"
	"golang.org/x/net/context"
)

func TestJoinRoomSuccess(t *testing.T) {
	l := LobbyMock{
		OnJoinRoom: func(roomID, userID string) error {
			return nil
		},
	}
	s := service.JoinRoom{
		Service: &service.GoRoom{
			Lobby: l,
		},
	}

	ps := httpservice.EmptyParams()
	ps.Set("roomID", "room1")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader(`{"user_id": "user123"}`))
	ctx := httpservice.WithParams(context.Background(), ps)
	err := s.DoHTTP(ctx, w, req)
	if err != nil {
		t.Fatalf("Joining a room should no fail but got: %s", err)
	}
	if w.Code != http.StatusOK {
		t.Errorf("Should respond with status 'OK' but got: %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Unexpected content type: %s", w.Header().Get("Content-Type"))
	}
	res := make(map[string]interface{})
	expected := map[string]interface{}{
		"status": "joined",
	}
	dec := json.NewDecoder(w.Body)
	if err := dec.Decode(&res); err != nil || !reflect.DeepEqual(res, expected) {
		t.Errorf("Invalid response body: expected '%v' but got '%v'", expected, res)
		if err != nil {
			t.Logf("Error: %s", err)
		}
	}
}

func TestJoinRoomMissingUserID(t *testing.T) {
	s := service.JoinRoom{
		Service: &service.GoRoom{
			Lobby: LobbyMock{},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader(`{}`))
	err := s.DoHTTP(context.Background(), w, req)
	if err == nil {
		t.Fatal("Expecting error but got none")
	}
}

func TestJoinRoomBadRequestBody(t *testing.T) {
	s := service.JoinRoom{
		Service: &service.GoRoom{
			Lobby: LobbyMock{},
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", strings.NewReader(`{`))
	err := s.DoHTTP(context.Background(), w, req)
	if err == nil {
		t.Fatal("Expecting error but got none")
	}
}
