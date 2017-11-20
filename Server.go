package dogo

import (
	"net/http"
	"strings"
	"time"
)

const(
	DEFAULT_ADDR = "127.0.0.1:0112"
	DEFAULT_READ_TIMEOUT = 1 * time.Second
	DEFAULT__WRITE_TIMEOUT = 30 * time.Second
)
type Server struct {
	ser *http.Server
	dogo *Dogo
}

type ServerConfig struct {
	Addr string `json:"addr"`
	ReadTimeout int `json:"readTimeout"`
	WriteTimeout int `json:"writeTimeout"`
	KeepaliveTimeout int `json:"keepaliveTimeout"`
}




func NewServer(dogo *Dogo) (*Server, error) {
	ser := &Server{
		ser : &http.Server{
			Addr : DEFAULT_ADDR,
			ReadTimeout : DEFAULT_READ_TIMEOUT,
			WriteTimeout : DEFAULT__WRITE_TIMEOUT,
		},
		dogo : dogo,
	}
	ser.init()
	return ser, nil
}


func(s *Server) init() {
	var sc ServerConfig
	err := s.dogo.Config.Get("server", &sc)
	if err != nil {
		s.dogo.Error("config ")
	}
	if !strings.EqualFold(sc.Addr, "") {
		s.ser.Addr = sc.Addr
	}

	if sc.ReadTimeout != 0 {
		s.ser.ReadTimeout = time.Duration(sc.ReadTimeout) * time.Second
	}

	if sc.WriteTimeout != 0 {
		s.ser.WriteTimeout = time.Duration(sc.WriteTimeout) * time.Second
	}

	/*if sc.KeepaliveTimeout != 0 {
		s.ser.IdleTimeout = time.Duration(sc.KeepaliveTimeout) * time.Second
	}*/

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		s.handleRequest(w, r)
	})
	s.ser.Handler = mux
}

func(s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	//构造context
	//初始化相关数据结构
	//构造request
	//启动链式调度
}

//This method will block the goroutine, please call it in a new goroutine if you don't want to be blocked
func(s *Server)Start() {
	err := s.ser.ListenAndServe()
	if err != nil {

	}
	//something big happen, server exit
}

