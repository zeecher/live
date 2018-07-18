package socket

import (
  "net/http"
  "log"
  "github.com/gorilla/websocket"
  "github.com/zeecher/live/store"
  "github.com/zeecher/live/utils"
  "encoding/json"
  "strconv"
)


var upgrade = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func init() {


    //register user commands
    CommandHandler.RegisterCommand("unload", func(messageJson map[string]interface{},user *store.User, uStore *store.Store)(bool) {
        log.Println("command is unload")
        log.Printf("users slice length before unload: %d", uStore.CountUsers())
        user.GetConn().Close()
        uStore.RemoveUserById(user.GetId())
        log.Printf("users slice length after unload: %d", uStore.CountUsers())
        return true
    })

    //register user commands
    CommandHandler.RegisterCommand("additional", func(messageJson map[string]interface{}, user *store.User, uStore *store.Store)(bool) {

        if !user.GetReady() {
            return false
        }

        // command is additional
        // add event to additional map
        // send rates to client

        eventID, ok := messageJson["eventID"]

        if !ok {
            panic("Error: Area ID is missing from submitted data.")
        }

        log.Printf("command is additional and eventID is %v\n", eventID)

        // convert eventID interface to int
        iEventID := utils.InterfaceToInt(eventID)

        log.Printf("one step back %v\n", uStore.GetOneStepBackEventOdds(iEventID))

        if command := uStore.GetOneStepBackEventOdds(iEventID);  command != "" {

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

            if err := user.GetConn().WriteMessage(websocket.TextMessage, []byte(formatted)); err != nil {
                log.Printf("error on message delivery e: %s\n", err)
            } else {
                log.Printf("user %s found, message sent %v\n",user.GetId(), formatted)
            }

            // append event id to additional slice if not specified
            user.SetAdditional(utils.AppendToSliceIfMissing(user.GetAdditional(), iEventID))
        }

        return false
    })

    //register user commands
    CommandHandler.RegisterCommand("mainline", func(messageJson map[string]interface{},user *store.User, uStore *store.Store)(bool) {

        if !user.GetReady() {
            return false
        }

        // convert event id from type interface to int
        eventID, ok := messageJson["eventID"]

        if !ok {
            panic("Error: Area ID is missing from submitted data.")
        }

        // convert eventID interface to int
        iEventID := utils.InterfaceToInt(eventID)

        // remove eventId from additional
        user.SetAdditional(utils.RemoveFromSliceIfExists(user.GetAdditional(), iEventID))

        // send mainlines only to the client

        m := map[string]interface{}{}

        err := json.Unmarshal([]byte(uStore.GetOneStepBackEventOdds(iEventID)), &m)

        if err != nil {
            panic(err)
        }

        // get rid from the list of opened additional
        if !utils.Contains(user.GetAdditional(), iEventID) {
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

        if err := user.GetConn().WriteMessage(websocket.TextMessage, []byte(formatted)); err != nil {
            log.Printf("error on message delivery e: %s\n", err)
        } else {
            log.Printf("user %s found, message sent\n", m)
        }

        return false

    })

    //register user commands
    CommandHandler.RegisterCommand("events", func(messageJson map[string]interface{}, user *store.User, uStore *store.Store) (bool) {

        for _, record := range messageJson {
            if events, ok := record.([]interface{}); ok {
                for _, event := range events {

                    if eventIDstr, ok := event.(string); ok {
                        /* act on str */
                        eventIDint, err := strconv.Atoi(eventIDstr)
                        if err == nil {
                            user.SetLive(utils.AppendToSliceIfMissing(user.GetLive(), eventIDint))
                        }

                    } else {
                        /* not string */
                        panic(ok)
                    }

                }

            }
        }

        // set user's status to ready in order to receive rates
        user.SetReady(true)

        //log.Printf("user lives %v\n", user.GetLive())

        return false
    })

}


func Handler(w http.ResponseWriter, r *http.Request, uStore *store.Store) {

  conn, err := upgrade.Upgrade(w, r, nil)

  if err != nil {
    log.Fatalf("upgrader error %s\n" + err.Error())
    return
  }

  defer conn.Close()

  //log.Printf("user slice length before join: %v\n", uStore.GetUsers())

  //add new user to the store
  user := uStore.NewUser(conn)

  //log.Printf("user %s joined\n", user.GetId())

  //log.Printf("user slice length after join: %v\n", uStore.GetUsers())

  for {

      messageType, messageByte, err := user.GetConn().ReadMessage()

      if ce, ok := err.(*websocket.CloseError); ok {

          switch ce.Code {

          case websocket.CloseNormalClosure,
               websocket.CloseGoingAway,
               websocket.CloseUnsupportedData,
               websocket.CloseNoStatusReceived:

               user.GetConn().Close()
               uStore.RemoveUserById(user.GetId())

              return
          }
      }

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

      if command, ok := messageJson["command"].(string); ok {

          if exit := CommandHandler.ExecuteCommand(command, messageJson, user, uStore); exit {
              return
          }

      } else {
          log.Fatalln("could not ran command")
      }

    }

}
