package pusuber


import (
  "testing"
  "github.com/rafaeljusto/redigomock"
  "github.com/gomodule/redigo/redis"

)

func TestRedisPubSub(t *testing.T) {
  redisChannel := "example"

  c := redigomock.NewConn()

  c.Command("SUBSCRIBE", redisChannel).Expect([]interface{}{
    []byte("subscribe"),
    []byte(redisChannel),
    []byte("1"),
  })
  c.Command("publish", redisChannel, "Uno").Expect([]interface{}{
    []byte("message"),
    []byte(redisChannel),
    []byte("Uno"),
  })
  c.Command("publish", redisChannel, "Deux").Expect([]interface{}{
    []byte("message"),
    []byte(redisChannel),
    []byte("Deux"),
  })
  c.Command("publish", redisChannel, "Three").Expect([]interface{}{
    []byte("message"),
    []byte(redisChannel),
    []byte("Three"),
  })
  c.Command("publish", redisChannel, "Finish").Expect([]interface{}{
    []byte("message"),
    []byte(redisChannel),
    []byte("Finish"),
  })

  messages := ReadMessages(c, redisChannel)

  t.Log("Started publisher connection")
  SendMessages(c, redisChannel, "Uno")
  SendMessages(c, redisChannel, "Deux")
  SendMessages(c, redisChannel, "Three")
  SendMessages(c, redisChannel, "Finished")
  t.Log("Sent messages")

  counter := 0
  for msg := range messages {
    t.Logf("Recieved a message: %s ", msg)
    counter++
    if msg == "Finished" {
      break
    }
  }
  if counter != 3 {
    t.Logf("Failed, expected 3 messages, recieved %d", counter)
    t.Fatal("Incorrect message count")
    return
  }
  t.Log("Succesful redis pub sub test")
}

func SendMessages(c redis.Conn, redisChannel, value string) {
  c.Send("publish", redisChannel, value)
}

type FRedis struct {
  Network string
  Address string
}

func (r *FRedis) GetNewConn() redis.Conn {

  return redigomock.NewConn()
}



func TestHandleChannel(t *testing.T)  {

  r := &FRedis{}

  HandleChannel("live", r )

}
