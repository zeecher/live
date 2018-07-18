package main

import (
  "log"
  "fmt"
  "net/http"
  "github.com/zeecher/live/utils"
  "github.com/zeecher/live/store"
  "github.com/zeecher/live/socket"
  "encoding/json"
  "github.com/gorilla/websocket"
  "github.com/zeecher/live/pusuber"
)

var serverAddress = ":5005"

type Message struct {
  MsgStr string
  MsgMap map[string]interface{}
}

var rateChannel chan Message
var finishedChannel chan Message
var universalChannel chan Message

var uStore  *store.Store

//to inform users about commands that used for general uncased purpose
func universalInformer(universalChannel chan Message) {

  for msg := range universalChannel {

    fmt.Printf("got universal command %v\n", msg)

    for _, u := range uStore.GetUsers() {

      if !u.GetReady() {
        continue
      }

      if err := u.GetConn().WriteMessage(websocket.TextMessage, []byte(msg.MsgStr)); err != nil {
        log.Printf("error on message delivery e: %s\n", err)
      } else {
        log.Printf("user found, message sent command %s\n", msg.MsgStr)
      }
    }
  }

}
func rateInformer(rateChannel chan Message) {

  for msg := range rateChannel {

    fmt.Printf("rate %v\n", msg)

    log.Printf("get users %v\n", uStore.GetUsers())

    eventID, ok := msg.MsgMap["eventID"]

    if !ok {
      log.Panic("Error: Area ID is missing from submitted data.")
      return
    }

    floatEventID, ok := eventID.(float64)

    iEventID := int(floatEventID)

    for _, u := range uStore.GetUsers() {


      if !u.GetReady() {
        continue
      }

      log.Printf("inside range %v\n", u.GetId())

      // if not in the lives add it to the list of lives
      // after adding inform user with new-rate command structure, change command rate to new-rate
      if !utils.Contains(u.GetLive(), iEventID) {
        u.SetLive(append(u.GetLive(), iEventID))
        msg.MsgMap["command"] = "new-rate"
      }

      if !utils.Contains(u.GetAdditional(), iEventID) {
        for _, record := range msg.MsgMap {
          if odds, ok := record.(map[string]interface{}); ok {

            delete(odds, "additional")
          }
        }
      }

      formatted, err := json.Marshal(msg.MsgMap)

      if err != nil {
        panic(err)
      }

      if err := u.GetConn().WriteMessage(websocket.TextMessage, []byte(formatted)); err != nil {
        log.Printf("error on message delivery e: %s\n", err)
      } else {
        log.Printf("user found, message sent command %s\n", formatted)
      }


      // add rate map to one step back dictionary
      uStore.AppendToOneStepBack(iEventID,  msg.MsgStr)

      log.Printf("rates one stap back %v\n", uStore.GetOneStepBackEventOdds(iEventID))

      // then live event is finished
      // remove it from the lives slice and additionals slice

    }

  }
}
func finishedInformer(finishedChannel chan Message) {

  for msg := range finishedChannel {
    fmt.Printf("finished %v\n", msg)
  }
}

func init() {
  //slice to store active users
  uStore = &store.Store{}
  uStore.InitUsers()
  uStore.InitOneStepBack()
}

func main() {

  rateChannel = make(chan Message)

  finishedChannel = make(chan Message)

  universalChannel = make(chan Message)

  go universalInformer(universalChannel)

  go rateInformer(rateChannel)

  go finishedInformer(finishedChannel)

  go processChannel("live")

  http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    socket.Handler(w, r, uStore)
  })

  log.Printf("server started at %s\n", serverAddress)

  log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func processChannel(channel string) {

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

func handleInformers(msgStr string, msgMap map[string]interface{}) {

  message := Message{msgStr,msgMap }

  // check which command is received and perform actions according to the command type
  switch  command := msgMap["command"]; command {

  case "rate":

    rateChannel <-message

  case "finished":

    finishedChannel <- message

  default:

    universalChannel <- message

  }

}





