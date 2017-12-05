package dogo

import (
	"testing"
	"reflect"
	"fmt"
)

type typeRegisStruct struct {
	Ip string `Value:"127.0.0.1"`
	Auth Authrization `Ref:"Authrization"`
	Pro Protocal `Autowired:"true"`
}

func(tr typeRegisStruct)Init() {
	fmt.Println("typeRegisStruct init")
}
type Authrization struct {
	Username string `Value:"dzthink"`
	Password string `Value:"123456"`
	Tokens []string
}

func(a Authrization)Init() {
	fmt.Println("hello")
}

type Protocal interface {
	Post()
}

type HttpProtocal struct {
}

func(hp *HttpProtocal)Post() {

}


func TestNewCtx(t *testing.T) {
	ctx := NewCtx([]*TypeMeta{})
	if ctx == nil {
		t.Error("fail to create ctx")
	}
}

func TestCtx_GetInstanceWithId(t *testing.T) {
	conf, err := NewConfig("config.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	ctx := NewCtx([]*TypeMeta{
		&TypeMeta{"dogo/Protocal", reflect.TypeOf((*Protocal)(nil)).Elem(), reflect.TypeOf(&HttpProtocal{})},
		&TypeMeta{"",reflect.TypeOf(typeRegisStruct{}), reflect.TypeOf(typeRegisStruct{})},
		&TypeMeta{"dogo/authrization",reflect.TypeOf(Authrization{}), reflect.TypeOf(Authrization{})},
	})
	ctxConf, err := conf.ChildList(CONF_CTX)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	ctx.parseBluePrint(ctxConf)
	tmp := ctx.GetInstanceWithType(reflect.TypeOf(typeRegisStruct{}))
	fmt.Println(tmp)
	ctx.active()
}
