package store

import (
  "sync"
  "os/exec"
  "github.com/gorilla/websocket"
  "log"
)


type Store struct {
  Users []*User
  OneStepBack map[int]string
  sync.Mutex

}

func (s *Store) NewUser(conn *websocket.Conn) *User {

  out, err := exec.Command("uuidgen").Output()
  if err != nil {
    log.Fatal(err)
  }

  u := &User{
    ID:   string(out),
    conn: conn,
    ready: false,
  }
  s.Lock()
  defer s.Unlock()
  s.Users = append(s.Users, u)
  return u
}



