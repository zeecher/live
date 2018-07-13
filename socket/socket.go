package socket

import (
  "net/http"
  "log"
  "github.com/gorilla/websocket"
  "github.com/zeecher/live/store"
  "encoding/json"
)


var upgrade = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func Handler(w http.ResponseWriter, r *http.Request, uStore *store.Store) {

  conn, err := upgrade.Upgrade(w, r, nil)

  if err != nil {
    log.Fatalf("upgrader error %s\n" + err.Error())
    return
  }

  defer conn.Close()

  log.Printf("user slice length before join: %v\n", uStore.GetUsers())

  //add new user to the store
  user := uStore.NewUser(conn)

  log.Printf("user %s joined\n", user.GetId())

  log.Printf("user slice length after join: %v\n", uStore.GetUsers())


  for {

      messageType, messageByte, err := user.GetConn().ReadMessage()

      if err != nil {
        log.Printf("socket connection read message error %s\n",err.Error())
        continue
      }

      // if message is type of string go inside code block
      // if not ignore because don't know how to handle json for now
      if messageType != 1 {
        log.Printf("unknown message type recieved from clinet: %T\n", messageType)
        continue
      }

      messageStr := string(messageByte)
      log.Printf("recieved message from clinet: %s", messageStr)

      messageJson := map[string]interface{}{}
      if  err := json.Unmarshal([]byte(messageStr), &messageJson); err != nil {
        log.Printf("message string json unmarshalling error: %s\n", err.Error())
        continue
      }

      // check which command is received from client and perform actions according to the command type
      switch  command := messageJson["command"]; command {

        case "unload":
          log.Println("command is unload")
          log.Printf("users slice length before unload: %d", uStore.CountUsers())
          user.GetConn().Close()
          uStore.RemoveUserById(user.GetId())
          log.Printf("users slice length after unload: %d", uStore.CountUsers())
          return

      }
    }

}
