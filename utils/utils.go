package utils

import "os/exec"
import "log"


func Contains(s []int, e int) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

func RemoveFromSliceIfExists(l []int, item int) []int {
  for i, other := range l {
    if other == item {
      return append(l[:i], l[i+1:]...)
    }
  }
  return l
}


func InterfaceToInt(eventID interface{}) int  {
  floatEventID, _ := eventID.(float64)
  return int(floatEventID)
}

func AppendToSliceIfMissing(slice []int, i int) []int {
  for _, ele := range slice {
    if ele == i {
      return slice
    }
  }
  return append(slice, i)
}

func GenUUID() string {

  out, err := exec.Command("uuidgen").Output()

  if err != nil {
    log.Fatal(err)
  }

  return string(out)

}
