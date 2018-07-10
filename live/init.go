package live

import (

  "github.com/zeecher/live/store"
)


var UserStore  *store.Store
var PubSubStore  *PubSub


func init() {

  //slice to store active users
  UserStore  = &store.Store{
    Users: make([]*store.User, 0, 1),
    OneStepBack: map[int]string{},

  }

  PubSubStore = &PubSub{}

  PubSubStore.SetRedisConn()
  PubSubStore.SetPubSubConn()

}


