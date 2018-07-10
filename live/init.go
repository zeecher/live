package live

import (

  "github.com/zeecher/live/store"
)


var UserStore  *store.Store


func init() {

  //slice to store active users
  UserStore  = &store.Store{
    Users: make([]*store.User, 0, 1),
    OneStepBack: map[int]string{},
  }

}


