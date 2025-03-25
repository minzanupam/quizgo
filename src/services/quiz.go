package services

import (
	"fmt"
	"log"
	"net/http"
	"quizgo/src/views"
)

func (s *Service) quizParentPageHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := authenticate(s.Store, r); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login?redirect_url=%2Fquiz", http.StatusTemporaryRedirect)
		return
	}
	page := views.QuizParentPage()
	if err := page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) quizApiHandler(w http.ResponseWriter, r *http.Request) {
	quizTitle := r.FormValue("title")
	row := s.Db.QueryRow(r.Context(), `INSERT INTO (title) VALUES ($1) RETURNING ID`, quizTitle)
	var quizID int64
	if err := row.Scan(&quizID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/dashboard/quiz/%ld", quizID), http.StatusTemporaryRedirect)
}
