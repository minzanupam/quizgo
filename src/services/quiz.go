package services

import (
	"fmt"
	"log"
	"net/http"
	"quizgo/src/views"
	"strconv"
	"time"
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
	userID, err := authenticate(s.Store, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	quizTitle := r.FormValue("quiz_title")
	if quizTitle == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row := s.Db.QueryRow(r.Context(), `INSERT INTO quizzes (title, owner_id) VALUES ($1, $2) RETURNING ID`, quizTitle, userID)
	var quizID int64
	if err := row.Scan(&quizID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/dashboard/quiz/%d", int(quizID)), http.StatusFound)
}

func (s *Service) quizPageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(s.Store, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	quizID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row := s.Db.QueryRow(r.Context(), `SELECT title, created_at, updated_at
		FROM quizzes WHERE ID = $1 AND owner_id = $2`, quizID, userID)

	var quiz views.DBQuiz
	var createdAt time.Time
	var updatedAt time.Time
	if err = row.Scan(&quiz.Title, &createdAt, &updatedAt); err != nil {
		log.Println(err)
	}
	quiz.ID = strconv.Itoa(quizID)
	quiz.CreatedAt = createdAt.Format(time.RFC3339)
	quiz.UpdatedAt = updatedAt.Format(time.RFC3339)

	page := views.QuizPage(quiz)
	if err := page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
