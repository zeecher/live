package pusuber

import (
  "github.com/gomodule/redigo/redis"
)

type IRedis interface {

  GetNewConn() redis.Conn
}

type Redis struct {
  Network string
  Address string
}

func (r *Redis) GetNewConn() redis.Conn {

  conn , err := redis.Dial(r.Network, r.Address)

  if err != nil {
    panic(err)
  }

  return conn
}


