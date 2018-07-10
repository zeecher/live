package store

import (
  "github.com/gorilla/websocket"
)

type IUser interface {

  SetReady(r bool)
  GetReady()
  GetLive()
  SetLive(lives []int)
  SetAdditional(additional []int)
  GetAdditional()
  GetConn()
}

type User struct {
  ID   string
  conn *websocket.Conn
  additional []int // to hold list of additional
  live []int // to hold list of lives
  ready bool
}

func (s *User) SetReady(r bool)  {
  s.ready = r
}

func (s *User) GetReady() bool {
  return s.ready
}

func (s *User) GetLive() []int {
  return s.live
}

func (s *User) SetLive(lives []int) {
  s.live = lives
}

func (s *User) SetAdditional(additional []int) {
  s.additional = additional
}

func (s *User) GetAdditional() []int {
 return s.additional
}

func (s *User) GetConn() *websocket.Conn {
  return s.conn
}