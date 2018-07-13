package store

import (
  "sync"
  "github.com/gorilla/websocket"
  "github.com/zeecher/live/utils"
  "log"
)


type Store struct {
  users []*User
  oneStepBack map[int]string
  sync.Mutex
}
func (s *Store) GetUsers() []*User{
  return s.users
}

func (s *Store) InitUsers() {
  s.users =  make([]*User, 0, 1)
}

func (s *Store) NewUser(conn *websocket.Conn) *User {

  u := &User{
    id:   utils.GenUUID(utils.RealCommander{}),
    conn: conn,
    ready: false,
  }
  s.Lock()
  defer s.Unlock()
  s.users = append(s.users, u)
  return u
}

// Iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builtin range keyword
func (s *Store) Iterator() <-chan User {
  c := make(chan User)

  f := func() {
    s.Lock()
    defer s.Unlock()
    for _, value := range s.users {
      c <- User{value.GetId(), value.GetConn(), value.GetAdditional(), value.GetLive(), value.GetReady()}
    }
    close(c)
  }
  go f()

  return c
}

func (s *Store) GetOneStepBackEventOdds(eventId int) string {

  s.Lock()
  defer s.Unlock()

  if odds, ok :=  s.oneStepBack[eventId]; ok {

    return odds
  }

  return ""
}

func (s *Store) AppendToOneStepBack(eventId int, odds string) {

  s.Lock()
  defer s.Unlock()
  s.oneStepBack[eventId] = odds
}

func (s *Store) CountUsers() int {

  return len(s.users)
}

func (s *Store) RemoveUserById(userId string) {

  s.Lock()
  defer s.Unlock()

  for index, user := range s.users {

    if user.GetId() == userId {

      log.Printf("removing user %s\n", userId)

      s.users = append(s.users[:index], s.users[index+1:]...)
    }
  }

}