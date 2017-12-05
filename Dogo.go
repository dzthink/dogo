package dogo

import (
	"reflect"
	"syscall"
	"os"
	"sync"
)

var(
	config *Config
	context *Ctx
	sigs chan os.Signal
	wg sync.WaitGroup
)

const(
	CONF_CTX = "dogo.ctx"
	CONF_CTX_ID = "id"
	CONF_CTX_INIT_METHOD = "init-method"
	CONF_CTX_ALIAS = "alias"
	CONF_CTX_FIELDS = "fields"
	CONF_CTX_FIELD_NAME = "name"
	CONF_CTX_FIELD_TYPE = "type"
	CONF_CTX_FIELD_VALUE = "value"
)


func GetConfig() *Config {
	return config
}

func NewContext(ts []*TypeMeta, confPath string) {
	//process panic, trigger signal SIGUSER2
	defer func() {
		if err := recover(); err != nil {
			sigs <- syscall.SIGINT
		}
	}()

	var err error
	//初始化配置
	config, err = NewConfig(confPath)
	if err != nil {
		panic("Can not create config:" + err.Error())
	}
	//todo 准备环境：进程用户/组,pid记录
	logProcessInfo(config)
	//信号处理
	wg.Add(1)
	go func() {
		processSignal()
		wg.Done()
	}()

	//context初始化及处理
	wg.Add(1)
	ts = ensureLogImplementExist(ts)
	context = NewCtx(ts)
	ctxConfig, err := config.ChildList(CONF_CTX)
	if err != nil {
		panic("Context config error:" + err.Error())
	}
	context.parseBluePrint(ctxConfig)

	go func() {
		context.active()
		wg.Done()
	}()
	wg.Wait()
}

func ensureLogImplementExist(ts []*TypeMeta) []*TypeMeta{
	logIf := reflect.TypeOf((*Log)(nil)).Elem()
	for _, tm := range ts {
		if tm.Abstract == logIf {
			return ts
		}
	}
	return append(ts, &TypeMeta{"", logIf, reflect.TypeOf(&DogoLog{})})
}

func GetContext() *Ctx {
	return context
}

func activeContext() {

}

