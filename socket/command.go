package socket

import (
  "errors"
  "github.com/zeecher/live/store"

)

/*
 * How to avoid a long switch-case statement
 * https://stackoverflow.com/questions/44812324/how-to-avoid-a-long-switch-case-statement-in-go
 */
type HandlerFunc func(commandMap map[string]interface{},user *store.User, uStore *store.Store)(bool)
type HandlerMap map[string]HandlerFunc
var  CommandHandler = HandlerMap{}


func (hr HandlerMap) RegisterCommand (command string, f HandlerFunc) error {

  if _, exists := hr[command]; exists {

    return errors.New("command already exists")
  }

  hr[command] = f

  return nil
}

func (hr HandlerMap) ExecuteCommand(command string, commandMap map[string]interface{}, user *store.User, uStore *store.Store) bool {

  if com, exists := hr[command]; exists {

    return com(commandMap,user,uStore)
  }

  return false
}
