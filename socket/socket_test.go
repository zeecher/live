package socket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gorilla/websocket"
	"github.com/zeecher/live/store"
)

var uStore  *store.Store
var first string
var second string
var third string

func init() {

	first := `{"mainlines":{"1":"1", "2": "2", "3": "3"},"additional": {"1":"1", "2": "2", "3": "3"}}`
	second := `{"mainlines":{"1":"1", "2": "2", "3": "3"},"additional": {"1":"1", "2": "2", "3": "3"}}`
	third := `{"mainlines":{"1":"1", "2": "2", "3": "3"},"additional": {"1":"1", "2": "2", "3": "3"}}`

	//slice to store active users
	uStore = &store.Store{}
	uStore.InitUsers()
	uStore.InitOneStepBack()

	uStore.AppendToOneStepBack(1, first)
	uStore.AppendToOneStepBack(2, second)
	uStore.AppendToOneStepBack(3, third)
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}


func TestSocket(t *testing.T) {

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	for i := 0; i < 10; i++ {
		if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
			t.Fatalf("%v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if string(p) != "hello" {
			t.Fatalf("bad message")
		}
	}
}


func TestHandler(t *testing.T) {

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		Handler(w, r, uStore)
	}))

	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"command": "events", "events": ["1", "2"]}`)); err != nil {
		t.Fatalf("%v", err)
	}

	if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"command": "additional", "eventID": 1 }`)); err != nil {
		t.Fatalf("%v", err)
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("from socket to client: %s", p)

	if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"command": "mainline", "eventID": 1 }`)); err != nil {
		t.Fatalf("%v", err)
	}

	/*_, p, err = ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("from socket to client: %s", p)
*/
/*	if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"command": "unload"}`)); err != nil {
		t.Fatalf("%v", err)
	}*/


}