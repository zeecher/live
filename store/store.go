package store

import (
  "sync"
  "github.com/gorilla/websocket"
  "github.com/zeecher/live/utils"
)

type IStore interface {

  NewUser(conn *websocket.Conn) *User
}

type Store struct {
  Users []*User
  OneStepBack map[int]string
  sync.Mutex

}


func (s *Store) NewUser(conn *websocket.Conn) *User {

  u := &User{
    ID:   utils.GenUUID(),
    conn: conn,
    ready: false,
  }
  s.Lock()
  defer s.Unlock()
  s.Users = append(s.Users, u)
  return u
}



