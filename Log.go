package dogo


type Log interface {
	Init(conf *Config)
	Debug(s string, v ...interface{})
	Info(s string, v ...interface{})
	Error(s string, v ...interface{})
	Fatal(s string, v ...interface{})
}

type DogoLog struct {

}

func(dg *DogoLog) Init(conf *Config) {

}

func(dg *DogoLog) Debug(s string, v ...interface{}) {

}

func(dg *DogoLog) Info(s string, v ...interface{}) {

}

func(dg *DogoLog) Error(s string, v ...interface{}) {

}

func(dg *DogoLog) Fatal(s string, v ...interface{}) {

}

