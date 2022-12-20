package internal

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func (s *Server) InitServer() {

	s.router = mux.NewRouter()
	s.viewrender = render.New(render.Options{
		IsDevelopment: true,
		Layout:        "layout",
		UnEscapeHTML:  true,
	})

	s.router.PathPrefix("/assets/").Handler(stat1)
	s.router.HandleFunc("/Tokens", s.ReadTokens).Methods("GET")
	s.router.HandleFunc("/", index)
	s.router.HandleFunc("/aboutUs", aboutUs)
	s.router.HandleFunc("/attorneys", attorneys)
	s.router.HandleFunc("/realEstate", realEstate)
	s.router.HandleFunc("/business", business)
	s.router.HandleFunc("/civilLitigation", civilLitigation)
	s.router.HandleFunc("/criminalLitigation", criminalLitigation)
	s.router.HandleFunc("/lawFirms", lawFirms)
	s.router.HandleFunc("/resources", resources)
	s.router.HandleFunc("/contactUs", contactUs)
	s.router.HandleFunc("/api/email", emailForm)
}
