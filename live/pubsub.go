package live

import (
  "github.com/garyburd/redigo/redis"
)

type PubSub struct {
  RedisConn redis.Conn
  PubSubConn *redis.PubSubConn
}

func (s *PubSub) SetPubSubConn()  {

  s.PubSubConn = &redis.PubSubConn{Conn: s.GetRedisConn()}
}

func (s *PubSub) GetPubSubConn() *redis.PubSubConn {

  return s.PubSubConn
}

func (s *PubSub) SetRedisConn() {

  conn , err := redis.Dial("tcp", ":6379")

  if err != nil {
    panic(err)
  }

  s.RedisConn = conn

}

func (s *PubSub) GetRedisConn() (redis.Conn) {

  return s.RedisConn
}

