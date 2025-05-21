package services

import (
	"log"
	"net/http"

	"quizgo/src/views"
)

func (s *Service) indexPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.HomePage().Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
