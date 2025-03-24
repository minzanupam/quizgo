package services

import (
	"log"
	"net/http"
)

func (s *Service) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.Store.Get(r, "authsession")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		log.Println(err)
	}
}
