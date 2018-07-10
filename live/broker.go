package live

import (
  "log"
  "encoding/json"
  "github.com/garyburd/redigo/redis"
  "github.com/gorilla/websocket"
  "github.com/zeecher/live/utils"
)

var PubSubConn *redis.PubSubConn

var RedisConn = func() (redis.Conn, error) {
  return redis.Dial("tcp", ":6379")
}

//pub/sub channel message broker
func MessageBroker()  {

  //open redis connection
  RedisConn , err := RedisConn()

  if err != nil {
    panic(err)
  }

  defer RedisConn.Close()

  //subscribe to pub/sub live channel
  PubSubConn = &redis.PubSubConn{Conn: RedisConn}

  PubSubConn.Subscribe("live")

  defer PubSubConn.Close()

  for {

    switch v := PubSubConn.Receive().(type) {
    case redis.Message:

      inform(v.Channel, string(v.Data))

      log.Printf("message received: %s: %s \n", v.Channel, string(v.Data))

    case redis.Subscription:
      log.Printf("subscription message: %s: %s %d\n", v.Channel, v.Kind, v.Count)

    case error:
      log.Fatal( v )
      return
    }
  }
}

func inform(channel string, content string)  {

  log.Println(content)

  m := map[string]interface{}{}

  // individual section for every users preferences
  for _, u := range UserStore.Users {

    if !u.GetReady() {
      continue
    }

    err := json.Unmarshal([]byte(content), &m)

    // if error in unmarshal panic about that
    if err != nil {
      panic(err)
    }

    // check which command is received and perform actions according to the command type
    switch  command := m["command"]; command {
    case "rate":

      // if command is rate
      // check that it is open or closed command
      // if open send all rates, if not just send mainlines
      log.Printf("rate command %v\n", m)

      eventID, ok := m["eventID"]

      if !ok {
        log.Panic("Error: Area ID is missing from submitted data.")
        return
      }

      floatEventID, ok := eventID.(float64)

      iEventID := int(floatEventID)

      // if not in the lives add it to the list of lives
      // after adding inform user with new-rate command structure, change command rate to new-rate
      if !utils.Contains(u.GetLive(), iEventID) {
        u.SetLive(append(u.GetLive(), iEventID))
        m["command"] = "new-rate"
      }

      if !utils.Contains(u.GetAdditional(), iEventID) {
        for _, record := range m {
          if odds, ok := record.(map[string]interface{}); ok {
            delete(odds, "additional")
          }
        }
      }


      formatted, err:= json.Marshal(m)

      if err != nil {
        panic(err)
      }

      if err := u.GetConn().WriteMessage(websocket.TextMessage, []byte(formatted)); err != nil {
        log.Printf("error on message delivery e: %s\n", err)
      } else {
        log.Printf("user %s found, message sent command %s\n", channel, formatted)
      }

      // add rate map to one step back dictionary
      UserStore.OneStepBack[iEventID] = content

      log.Printf("Added to %v",  UserStore.OneStepBack[iEventID])

      // then live event is finished
      // remove it from the lives slice and additional slice

    case "finished":

      // convert eventID from interface to int
      eventID, ok := m["eventID"]

      if !ok {
        log.Panic("Error: Area ID is missing from submitted data.")
        return
      }

      floatEventID, ok := eventID.(float64)

      iEventID := int(floatEventID)

      // if lives has this match remove it from lives
      u.SetLive(utils.RemoveFromSliceIfExists(u.GetLive(), iEventID))
      // remove eventId from additional
      u.SetAdditional(utils.RemoveFromSliceIfExists(u.GetAdditional(), iEventID))

      log.Printf("finished command %v\n", m)

      if err := u.GetConn().WriteMessage(websocket.TextMessage, []byte(content)); err != nil {
        log.Printf("error on message delivery e: %s\n", err)
      } else {
        log.Printf("user %s found, message sent\n", channel)
      }

    default:
      log.Printf("switch default %s\n", content)

      if err := u.GetConn().WriteMessage(websocket.TextMessage, []byte(content)); err != nil {
        log.Printf("error on message delivery e: %s\n", err)
      } else {
        log.Printf("user %s found, message sent\n", channel)
      }
    }

  }
}
