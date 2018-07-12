package utils

import "testing"

type TestCommander struct{

  UUID string
}

func (c TestCommander) Output(command string, args ...string) ([]byte, error) {
  return []byte(c.UUID) , nil
}

func TestGenUUID(t *testing.T)  {

  expected := "e3d38fc6-8d10-456d-bfb4-aa3a450f65c7"

  result := GenUUID(TestCommander{expected})

  if result != expected {
    t.Fatalf("Expected %s, got %s", expected, result)
  }
}
