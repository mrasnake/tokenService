package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"net/url"
)

func NewServer(srvc *TokenService) *Server {
	return &Server{
		service: srvc,
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

	ret, err := s.service.ReadTokens(req)
	if err != nil {
		return
	}

	s.viewrender.JSON(w, http.StatusOK, ret.tokenSecrets)

}

func (s *Server) WriteToken(w http.ResponseWriter, r *http.Request) {
	var req *WriteTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, err := s.service.WriteToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.viewrender.JSON(w, http.StatusOK, ret)
}

func (s *Server) UpdateToken(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var in TokenSecret

	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	in.Token = vars["token"]

	req := &UpdateTokenRequest{
		tokenSecret: in,
	}

	err = s.service.UpdateToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.viewrender.JSON(w, http.StatusNoContent, nil)
}

func (s *Server) DeleteToken(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	req := &DeleteTokenRequest{
		Token: vars["token"],
	}

	err := s.service.DeleteToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.viewrender.JSON(w, http.StatusNoContent, nil)

}
