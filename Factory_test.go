package dogo

import (
	"testing"
	"fmt"
)

type typeRegisStruct struct {
	Ip string `Value:"127.0.0.1"`
	Auth Authrization `Autowired:"dogo/authrization"`
}

type Authrization struct {
	Username string `Value:"dzthink"`
	Password string `Value:"123456"`
}

func TestFactory_RegisType(t *testing.T) {
	fac := NewFactory()
	fac.RegisType(typeRegisStruct{}, "dogo/typeRegisStruct", SCOPE_STATELESS)
	if _, ok := fac.typeMap["dogo/typeRegisStruct"]; !ok {
		t.Error("type register fail")
	}
}

func TestFactory_Make(t *testing.T) {
	fac := NewFactory()
	fac.RegisType(typeRegisStruct{}, "dogo/typeRegisStruct", SCOPE_STATELESS)
	fac.RegisType(Authrization{}, "dogo/authrization", SCOPE_STATELESS)
	if ins, err := fac.Make("dogo/typeRegisStruct"); err != nil {
		t.Error(err)
	} else {
		fmt.Println(ins)
		if makeResult, typeAssert := ins.(*typeRegisStruct); !typeAssert {
			t.Fatal("factory:make return wrong type")
		} else {
			if makeResult.Ip != "127.0.0.1" {
				t.Error("parameter injection fail")
			}
			if makeResult.Auth.Username != "dzthink" {
				t.Error("parameter injection fail")
			}
		}

	}
}
