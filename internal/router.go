package internal

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
)

func (s *Server) InitServer() {

	s.Router = mux.NewRouter()
	s.Viewrender = render.New(render.Options{
		IsDevelopment: true,
		Layout:        "layout",
		UnEscapeHTML:  true,
	})

	s.Router.HandleFunc("/tokens", s.ReadTokens).Methods(http.MethodGet)
	s.Router.HandleFunc("/tokens", s.WriteToken).Methods(http.MethodPost)
	s.Router.HandleFunc("/tokens/{token}", s.UpdateToken).Methods(http.MethodPut)
	s.Router.HandleFunc("/tokens/{token}", s.DeleteToken).Methods(http.MethodDelete)
}
