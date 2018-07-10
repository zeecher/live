package store

import (
  "github.com/gorilla/websocket"
)

type User struct {
  ID   string
  conn *websocket.Conn
  additional []int // to hold list of additional
  live []int // to hold list of lives
  ready bool
}

func (s *User) setReady(r bool)  {
  s.ready = r
}


func (s *User) getReady() bool {
  return s.ready
}
