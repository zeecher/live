package pusuber

import (
  "github.com/gomodule/redigo/redis"
)

func ReadMessages(c redis.Conn, redisChannel string) <-chan string {
  messages := make(chan string)

  go func() {
    psc := redis.PubSubConn{Conn: c}
    psc.Subscribe(redisChannel)
    for {
      switch v := psc.Receive().(type) {
      case redis.Message:
        messages <- string(v.Data)
      case error:
        //panic(v.Error())
        //messages <- "error:" + v.Error()
        close(messages)
        return
      }
    }
  }()

  return messages
}



/*func HandleChannels(channel string, r *IRedis)  {


  messages := ReadMessages(conn, channel)

  m := map[string]interface{}{}

  for msg := range messages {

    UnmarshalToInterface(msg,&m)

    //inform(m)
  }
}*/

func HandleChannel(c redis.Conn, redisChannel string) <-chan string {
  messages := make(chan string)

  go func() {
    psc := redis.PubSubConn{Conn: c}
    psc.Subscribe(redisChannel)
    for {
      switch v := psc.Receive().(type) {
      case redis.Message:
        messages <- string(v.Data)
      case error:
        close(messages)
        return
      }
    }
  }()

  return messages
}



