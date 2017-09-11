package dogo

import (
	"net/http"
)
type Server struct {
	ser *http.Server
	dogo *Dogo
}

type ServerConfig struct {
	Addr string `json:"addr"`


}
func NewServer(dogo *Dogo) (*Server, error) {
	ser := &Server{
		ser : &http.Server{},
		dogo : dogo,
	}
	return ser, nil
}


func(s *Server) init() error {
	conf, err := s.dogo.Config.Child("server")
	if err != nil {
		s.dogo.Error("config ")
	}
}

func(s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {

}

