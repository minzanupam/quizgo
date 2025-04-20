package services

import (
	"log"
	"net/http"
	"time"

	"quizgo/src/views"
)

func (s *Service) indexPageHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query(r.Context(), `
		SELECT ID, title, created_at, updated_at
		FROM quizzes
		WHERE status = 'published'`)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	quizzes := make([]views.DBQuiz, 0)
	for rows.Next() {
		var quiz views.DBQuiz
		var createdAt, updatedAt time.Time
		if err = rows.Scan(&quiz.ID, &quiz.Title, &createdAt, &updatedAt); err != nil {
			log.Println(err)
			continue
		}
		quiz.CreatedAt = createdAt.Format(time.RFC3339)
		quiz.UpdatedAt = updatedAt.Format(time.RFC3339)
		quizzes = append(quizzes, quiz)
	}
	if err = views.QuizListPage(quizzes).Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
