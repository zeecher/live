package main

import (

  "github.com/zeecher/live/pusuber"
  "github.com/zeecher/live/utils"
  "fmt"
  "sync"
  "github.com/zeecher/live/store"
)



type Message struct {
  MsgStr string
  MsgMap map[string]interface{}
}

var rateChannel chan Message
var finishedChannel chan Message

var UserStore  *store.Store

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
  UserStore = &store.Store{
    Users: make([]*store.User, 0, 1),
    OneStepBack: map[int]string{},

  }
}

func main() {

  wg := &sync.WaitGroup{}

  wg.Add(2)

  rateChannel = make(chan Message)

  finishedChannel = make(chan Message)

  go rateInformer(rateChannel)

  go finishedInformer(finishedChannel)

  go processChannel("live")

  wg.Wait()


 /* http.HandleFunc("/ws", live.WsHandler)

  log.Printf("server started at %s\n", live.ServerAddress)

  log.Fatal(http.ListenAndServe(live.ServerAddress, nil))*/
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





