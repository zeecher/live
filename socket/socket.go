package socket

import (
  "net/http"
  "log"
  "github.com/gorilla/websocket"
  )

var ServerAddress = ":5005"

var upgrade = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func Handler(w http.ResponseWriter, r *http.Request) {

  conn, err := upgrade.Upgrade(w, r, nil)
  if err != nil {
    log.Printf("upgrader error %s\n" + err.Error())
    return
  }

  UserStore.

}
