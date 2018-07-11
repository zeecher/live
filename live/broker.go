package live

import (
  "log"
  "encoding/json"
  "github.com/garyburd/redigo/redis"
  "github.com/gorilla/websocket"
  "github.com/zeecher/live/utils"
  "github.com/zeecher/live/store"
)


//pub/sub channel message broker
func MessageBroker()  {

  redisConn := PubSubStore.GetRedisConn()

  defer redisConn.Close()

  //subscribe to pub/sub live channel
  pubSubConn := PubSubStore.GetPubSubConn()

  pubSubConn.Subscribe("live")

  defer pubSubConn.Close()

  for {

    switch v := pubSubConn.Receive().(type) {
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


func rate(u *store.User, m map[string]interface{}, content string, channel string)  {

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

  sendToUser(u, formatted, channel)

  // add rate map to one step back dictionary
  UserStore.OneStepBack[iEventID] = content

  log.Printf("Added to %v",  UserStore.OneStepBack[iEventID])

  // then live event is finished
  // remove it from the lives slice and additional slice

}

func finished(u *store.User, m map[string]interface{},content string, channel string)  {

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

  sendToUser(u,[]byte(content),channel)
}

func sendToUser(u *store.User, formatted []byte, channel string)  {

  if err := u.GetConn().WriteMessage(websocket.TextMessage, []byte(formatted)); err != nil {
    log.Printf("error on message delivery e: %s\n", err)
  } else {
    log.Printf("user %s found, message sent command %s\n", channel, formatted)
  }
}

func unmarshalNested(content string) (map[string]interface{})  {

  var m = map[string]interface{}{}

  var err = json.Unmarshal([]byte(content), &m)

  // if error in unmarshal panic about that
  if err != nil {
    panic(err)
  }

  return m

}


func inform(channel string, content string)  {

  m := unmarshalNested(content)

  // individual section for every users preferences
  for _, u := range UserStore.Users {

    if !u.GetReady() {
      continue
    }

    // check which command is received and perform actions according to the command type
    switch  command := m["command"]; command {

    case "rate":

      rate(u, m, content, channel)

    case "finished":

      finished(u, m, content, channel)

    default:
      log.Printf("switch default %s\n", content)

      sendToUser(u, []byte(content), channel)
    }

  }
}
