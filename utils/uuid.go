package utils

import (
  "log"
  "os/exec"
  )

type Commander interface {
  Output(string, ...string) ([]byte, error)
}

type RealCommander struct{}

func (c RealCommander) Output(command string, args ...string) ([]byte, error) {
  return exec.Command(command, args...).Output()
}

func GenUUID(c Commander) string {

  out, err := c.Output("uuidgen")

  if err != nil {
    log.Fatal(err)
  }

  return string(out)
}
