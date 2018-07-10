package main

import (

  "github.com/zeecher/live/live"
  "net/http"
  "log"
)


func main() {

  go live.MessageBroker()

  http.HandleFunc("/ws", live.WsHandler)

  log.Printf("server started at %s\n", live.ServerAddress)

  log.Fatal(http.ListenAndServe(live.ServerAddress, nil))
}


