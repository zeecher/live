package main

import (

  "github.com/zeecher/live/pusuber"
  "github.com/zeecher/live/utils"
  "fmt"
  "github.com/zeecher/live/store"
  "github.com/zeecher/live/socket"
  "net/http"
  "log"
)

var serverAddress = ":5005"

type Message struct {
  MsgStr string
  MsgMap map[string]interface{}
}

var rateChannel chan Message
var finishedChannel chan Message

var uStore  *store.Store

func rateInformer(rateChannel chan Message)  {

  for msg := range rateChannel {
    fmt.Printf("rate %v\n", msg)
  }

}

func finishedInformer(finishedChannel chan Message)  {

  for msg := range finishedChannel {
    fmt.Printf("finished %v\n", msg)
  }
}

func init() {

  //slice to store active users
  uStore = &store.Store{}
  uStore.InitUsers()
}

func main() {

  rateChannel = make(chan Message)

  finishedChannel = make(chan Message)

  go rateInformer(rateChannel)

  go finishedInformer(finishedChannel)

  go processChannel("live")

  http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    socket.Handler(w, r, uStore)
  })

  log.Printf("server started at %s\n", serverAddress)

  log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func processChannel(channel string)  {

  redis := &pusuber.Redis{}
  redis.Network = "tcp"
  redis.Address = ":6379"

  conn :=  redis.GetNewConn()

  defer conn.Close()

  messages := pusuber.HandleChannel(conn, channel)

  msgMap := map[string]interface{}{}

  for msgStr := range messages {

    utils.UnmarshalToInterface(msgStr,&msgMap)

    handleInformers(msgStr,msgMap)
  }
}

func handleInformers(msgStr string, msgMap map[string]interface{})  {

  message := Message{msgStr,msgMap }

  // check which command is received and perform actions according to the command type
  switch  command := msgMap["command"]; command {

  case "rate":

    rateChannel <-message

  case "finished":

    finishedChannel <- message

  }

}





