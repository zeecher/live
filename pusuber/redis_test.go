package pusuber

import (
  "testing"
  "github.com/rafaeljusto/redigomock"
  "github.com/gomodule/redigo/redis"
  "fmt"
)

type FakeRedis struct {

}

func (r *FakeRedis) GetNewConn() *redigomock.Conn {

  return redigomock.NewConn()
}


type Person struct {
  Name string `redis:"name"`
  Age  int    `redis:"age"`
}

func RetrievePerson(conn redis.Conn, id string) (Person, error) {
  var person Person

  values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("person:%s", id)))
  if err != nil {
    return person, err
  }

  err = redis.ScanStruct(values, &person)
  return person, err
}

func TestRetrievePerson(t *testing.T) {

  fakeRedis := FakeRedis{}

  conn := fakeRedis.GetNewConn()

  cmd := conn.Command("HGETALL", "person:1").ExpectMap(map[string]string{
    "name": "Mr. Johson",
    "age":  "42",
  })

  person, err := RetrievePerson(conn, "1")
  if err != nil {
    t.Fatal(err)
  }

  if conn.Stats(cmd) != 1 {
    t.Fatal("Command was not called!")
  }

  if person.Name != "Mr. Johson" {
    t.Errorf("Invalid name. Expected 'Mr. Johson' and got '%s'", person.Name)
  }

  if person.Age != 42 {
    t.Errorf("Invalid age. Expected '42' and got '%d'", person.Age)
  }
}







