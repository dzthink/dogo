package dogo

import "reflect"

type Dogo struct {
	Factory *Factory
	Config *Config
	log Log
}

func NewDogo(path string, typeMap map[string]reflect.Type) *Dogo {
	conf, err := NewConfig(path)
	if err != nil {
		panic("config file parse fail")
	}
	//todo 初始化日志
	return &Dogo{
		Factory : NewFactory(typeMap),
		Config : conf,
		log : &DogoLog{},
	}
}


func(dg *Dogo)Info(s string, v ...interface{}) {

}

func(dg *Dogo)Debug(s string, v ...interface{}) {

}

func(dg *Dogo)Error(s string, v ...interface{}) {

}

func(dg *Dogo)Fatal(s string, v ...interface{}) {

}
