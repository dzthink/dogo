package dogo

import (
	"testing"
	"fmt"
)

type typeRegisStruct struct {
	Ip string `Value:"127.0.0.1"`
	Auth Authrization `Autowired:"dogo/authrization"`
	Pro Protocal `Autowired:"dogo/protocal"`
}

type Authrization struct {
	Username string `Value:"dzthink"`
	Password string `Value:"123456"`
}

type Protocal interface {
	Post()
}

type HttpProtocal struct {

}

func(hp *HttpProtocal)Post() {

}

func TestFactory_RegisType(t *testing.T) {
	typeMap := map[string]interface{} {
		"dogo/typeRegisStruct" : typeRegisStruct{},
	}
	fac := NewFactory(typeMap)
	if _, ok := fac.typeMap["dogo/typeRegisStruct"]; !ok {
		t.Error("type register fail")
	}
}

func TestFactory_Make(t *testing.T) {

	typeMap := map[string]interface{} {
		"dogo/typeRegisStruct" : typeRegisStruct{},
		"dogo/authrization" : Authrization{},
		"dogo/protocal" : HttpProtocal{},
	}
	fac := NewFactory(typeMap)

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
