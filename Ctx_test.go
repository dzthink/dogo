package dogo

import (
	"testing"
	"reflect"
	"fmt"
)

type typeRegisStruct struct {
	Ip string `Value:"127.0.0.1"`
	Auth Authrization `Autowired:"true"`
	Pro Protocal `Autowired:"true"`
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


func TestNewCtx(t *testing.T) {
	ctx := NewCtx()
	if ctx == nil {
		t.Error("fail to create ctx")
	}
}

func TestCtx_GetInstanceWithId(t *testing.T) {
	ctx := NewCtx()
	var p Protocal
	fmt.Println(reflect.TypeOf(&p).Elem())
	ctx.RegType(&TypeMeta{"dogo/Protocal", reflect.TypeOf(&p).Elem(), reflect.TypeOf(&HttpProtocal{})})
	ctx.RegType(&TypeMeta{"",reflect.TypeOf(typeRegisStruct{}), reflect.TypeOf(typeRegisStruct{})})
	ctx.RegType(&TypeMeta{"dogo/authrization",reflect.TypeOf(Authrization{}), reflect.TypeOf(Authrization{})})

	tmp := ctx.GetInstanceWithType(reflect.TypeOf(typeRegisStruct{}))
	fmt.Println(tmp)
}
