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
	ctx.RegType(&TypeMeta{"dogo/Protocal", reflect.TypeOf((*Protocal)(nil)).Elem(), reflect.TypeOf(&HttpProtocal{})})
	ctx.RegType(&TypeMeta{"",reflect.TypeOf(typeRegisStruct{}), reflect.TypeOf(typeRegisStruct{})})
	ctx.RegType(&TypeMeta{"dogo/authrization",reflect.TypeOf(Authrization{}), reflect.TypeOf(Authrization{})})

	bp := &Blueprint{
		TypeAlias : "dogo/authrization",
		Fields:make(map[string]*BluePrintField),
	}
	bp.Fields["Username"] = &BluePrintField{
		Name : "Username",
		Value : "zxduan",
		ValueType:ValueTypeConst,
	}
	bp.Fields["Tokens"] = &BluePrintField{
		Name : "Tokens",
		Value : []interface{}{"a", "b", "c"},
		ValueType : ValueTypeConst,
	}
	ctx.RegBlueprint("Authrization", bp)
	tmp := ctx.GetInstanceWithType(reflect.TypeOf(typeRegisStruct{}))
	fmt.Println(tmp)
}
