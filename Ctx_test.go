package dogo

import "testing"

type typeRegisStruct struct {
	Ip string `Value:"127.0.0.1"`
	Auth Authrization `Ref:"dogo/authrization"`
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
	//ctx.RegType()
}
