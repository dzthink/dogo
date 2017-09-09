package dogo

import (
	"testing"
	"fmt"
	"math"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestNewConfig(t *testing.T) {
	conf, err := NewConfig("config.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	timeout, err := conf.Int("timeout")
	if err != nil {
		t.Log(err)
	}
	if timeout != 1000 {
		t.Error("config parse error, wrong timeout")
	}

	var userlist []User
	conf.Get("users", &userlist)
	if len(userlist) != 2 || userlist[0].Username != "dzthink" {
		t.Log("struct parse fail")
		t.Fail()
	}

	vnd, err := conf.Float("vnd")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if math.Abs(vnd - 14.55) > 0.01 {
		fmt.Println(vnd)
		t.Log("float err")
		t.Fail()
	}
}