package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"net/url"
	"strings"
)

func NewServer(srvc *TokenService) *Server {
	return &Server{
		Service: srvc,
	}
}

type Server struct {
	Service    *TokenService
	Router     *mux.Router
	Viewrender *render.Render
}

func (s *Server) ReadTokens(w http.ResponseWriter, r *http.Request) {

	var toks []string
	m, _ := url.ParseQuery(r.URL.RawQuery)
	t := m["t"]

	if len(t) == 1 {
		toks = strings.Split(t[0], ",")
	} else {
		toks = t
	}

	req := &ReadTokenRequest{
		Tokens: toks,
	}

	ret, err := s.Service.ReadTokens(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Viewrender.JSON(w, http.StatusOK, ret.tokenSecrets)

}

func (s *Server) WriteToken(w http.ResponseWriter, r *http.Request) {
	var req *WriteTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, err := s.Service.WriteToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Viewrender.JSON(w, http.StatusOK, ret)
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

	err = s.Service.UpdateToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Viewrender.JSON(w, http.StatusNoContent, nil)
}

func (s *Server) DeleteToken(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	req := &DeleteTokenRequest{
		Token: vars["token"],
	}

	err := s.Service.DeleteToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Viewrender.JSON(w, http.StatusNoContent, nil)

}
