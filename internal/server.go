package internal

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"net/url"
)

func NewServer(srvc TokenService) *Server {
	return &Server{
		service: &srvc,
	}
}

type Server struct {
	service    *TokenService
	router     *mux.Router
	viewrender *render.Render
}

func (s *Server) ReadTokens(w http.ResponseWriter, r *http.Request) {

	m, _ := url.ParseQuery(r.URL.RawQuery)
	toks := m["t"]

	req := &ReadTokenRequest{
		Tokens: toks,
	}

	ret, err := s.service.ReadToken(req)
	if err != nil {
		return
	}

	s.viewrender.JSON(w, http.StatusOK, ret.tokenSecrets)

}

func (s *Server) WriteToken(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UpdateToken(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) DeleteToken(w http.ResponseWriter, r *http.Request) {

}
