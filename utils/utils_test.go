package utils

import (
  "testing"
  "reflect"
)

func TestAppendToSliceIfMissing(t *testing.T)  {

  want := []int{1,2,3,4,5,6}

  got := AppendToSliceIfMissing([]int{1,2,3,4,5}, 6)

  if !reflect.DeepEqual(want, got) {

    t.Fatalf("Want %v, got %v", want, got)

  }
}

func TestRemoveFromSliceIfExists(t *testing.T)  {

  want := []int{1,2,3,4,5,6}

  got := RemoveFromSliceIfExists([]int{1,2,3,4,5,6,7}, 7)

  if !reflect.DeepEqual(want, got) {

    t.Fatalf("Want %v, got %v", want, got)

  }

}

func TestContains(t *testing.T) {

  var got = true

  got = Contains([]int{1,2,3,4,5,6,7}, 4)

  if got != true {

    t.Fatalf("Want %v, got %v", true, got)

  }

  got = Contains([]int{1,2,3,4,5,6,7}, 9)

  if got != false {

    t.Fatalf("Want %v, got %v", false, got)
  }

}
