package internal

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer creates and returns a Server instance.
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

// ReadTokens servers as the transport layer GET function of /tokens endpoint,
// parsing the request, formatting the data and calling service layer function.
func (s *Server) ReadTokens(w http.ResponseWriter, r *http.Request) {

	// parse the token(s) from the URL Query String
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

	// call partner service layer function
	ret, err := s.Service.ReadTokens(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// write the response
	s.Viewrender.JSON(w, http.StatusOK, ret.tokenSecrets)

}

// WriteToken servers as the transport layer POST function of /tokens endpoint,
// parsing the request, formatting the data and calling service layer function.
func (s *Server) WriteToken(w http.ResponseWriter, r *http.Request) {
	var req *WriteTokenRequest

	// parse body of the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call partner service layer function
	ret, err := s.Service.WriteToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// write the response
	s.Viewrender.JSON(w, http.StatusOK, ret)
}

// UpdateToken servers as the transport layer PUT function of /tokens endpoint,
// parsing the request, formatting the data and calling service layer function.
func (s *Server) UpdateToken(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var in TokenSecret

	// parse body of the request
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	in.Token = vars["token"]

	req := &UpdateTokenRequest{
		tokenSecret: in,
	}

	// call partner service layer function
	err = s.Service.UpdateToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// write the response
	s.Viewrender.JSON(w, http.StatusNoContent, nil)
}

// DeleteToken servers as the transport layer DELETE function of /tokens endpoint,
// parsing the request, formatting the data and calling service layer function.
func (s *Server) DeleteToken(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	req := &DeleteTokenRequest{
		Token: vars["token"],
	}

	// call partner service layer function
	err := s.Service.DeleteToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// write the response
	s.Viewrender.JSON(w, http.StatusNoContent, nil)

}