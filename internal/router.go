package internal

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
)

func (s *Server) InitServer() {

	s.router = mux.NewRouter()
	s.viewrender = render.New(render.Options{
		IsDevelopment: true,
		Layout:        "layout",
		UnEscapeHTML:  true,
	})

	s.router.HandleFunc("/tokens", s.ReadTokens).Methods(http.MethodGet)
	s.router.HandleFunc("/tokens", s.WriteToken).Methods(http.MethodPost)
	s.router.HandleFunc("/tokens/{token}", s.UpdateToken).Methods(http.MethodPut)
	s.router.HandleFunc("/tokens/{token}", s.DeleteToken).Methods(http.MethodDelete)
