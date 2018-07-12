package live

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/websocket"
  "log"
  "strconv"
  "github.com/zeecher/live/utils"
)

var ServerAddress = ":5005"


var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func WsHandler(w http.ResponseWriter, r *http.Request) {

  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Printf("upgrader error %s\n" + err.Error())
    return
  }

  //add new user to the store
  u := UserStore.NewUser(conn)

  log.Printf("user %s joined\n", u.ID)

  for {

    messageType, messageByte, err := u.GetConn().ReadMessage()

    log.Printf("%+v\n", messageByte)

    if err != nil {
      log.Println(err)
      return
    }


    // if message is type of string go inside code block
    // if not ignore because don't know how to handle json for now
    if messageType == 1 {

      messageString := string(messageByte)

      m := map[string]interface{}{}

      err := json.Unmarshal([]byte(messageString), &m)

      // if error in unmarshal panic about that
      if err != nil {
        log.Println("fuck ci")
      }

      log.Println(m)

      // check which command is received from client and perform actions according to the command type
      switch  command := m["command"]; command {
      case "unload":

        log.Println("command is unload")

        log.Printf("users slice before %v", UserStore.Users)

        for index, users := range UserStore.Users {
          if users.ID == u.ID {
            u.GetConn().Close()

            UserStore.Users = append(UserStore.Users[:index], UserStore.Users[index+1:]...)

            log.Println("user removed")

            log.Printf("users slice after %v", UserStore.Users)

            return
          }
        }

      case "additional":
        // command is additional
        // add event to additional map
        // send rates to client


        if !u.GetReady() {
          continue
        }

        log.Printf("command is additional and eventID is %v\n", m["eventID"])

        // convert event id from type interface to int
        eventID, ok := m["eventID"]

        if !ok {
          panic("Error: Area ID is missing from submitted data.")
        }

        // convert eventID interface to int
        iEventID := utils.InterfaceToInt(eventID)

        log.Printf("here is one step back%v", UserStore.OneStepBack[iEventID])

        // check if oneStepBack contains eventID
        if command, ok := UserStore.OneStepBack[iEventID]; ok {

          // send mainlines only to the client

          m := map[string]interface{}{}

          err := json.Unmarshal([]byte(command), &m)

          if err != nil {
            panic(err)
          }

          m["command"] = "additional"

          formatted, err:= json.Marshal(m)

          if err != nil {
            panic(err)
          }

          if err := u.GetConn().WriteMessage(websocket.TextMessage, []byte(formatted)); err != nil {
            log.Printf("error on message delivery e: %s\n", err)
          } else {
            log.Printf("user found, message sent\n")
          }

          // append event id to additional slice if not specified
          u.SetAdditional(utils.AppendToSliceIfMissing(u.GetAdditional(), iEventID))

          log.Printf("user additional %v\n.", u.GetAdditional())
        }else {
          log.Printf("user additional is wrong %v\n.", u.GetAdditional())
        }

      case "mainline":
        // command to close additional for specific match

        if !u.GetReady() {
          continue
        }

        // convert event id from type interface to int
        eventID, ok := m["eventID"]

        if !ok {
          log.Panic("Error: Area ID is missing from submitted data.")
        }

        // convert eventID interface to int
        iEventID := utils.InterfaceToInt(eventID)

        // remove eventId from additional
        u.SetAdditional(utils.RemoveFromSliceIfExists(u.GetAdditional(), iEventID))

        // send mainlines only to the client

        m := map[string]interface{}{}

        err := json.Unmarshal([]byte(UserStore.OneStepBack[iEventID]), &m)

        if err != nil {
          panic(err)
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
          log.Printf("user %s found, message sent\n", m)
        }

      case "events":

        log.Printf("%s.", m)

        for _, record := range m {
          if events, ok := record.([]interface{}); ok {
            for _, event := range events {

              if eventIDstr, ok := event.(string); ok {
                /* act on str */
                eventIDint, err := strconv.Atoi(eventIDstr)
                if err == nil {
                  u.SetLive(utils.AppendToSliceIfMissing(u.GetLive(), eventIDint))
                }

              } else {
                /* not string */
                panic(ok)
              }

            }

          }
        }

        // set user's status to ready in order to receive rates
        u.SetReady(true)

        log.Printf("user lives %v\n", u.GetLive())

      default:
        log.Printf("%s.", command)
      }

    }
  }
}



