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

type QuizRow struct {
	ID           int64
	Title        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	QuestionID   int64
	QuestionBody string
}

func parseRowsToQuiz(rows []QuizRow) views.DBQuiz {
	var quiz views.DBQuiz
	for _, row := range rows {
		quiz.ID = strconv.Itoa(int(row.ID))
		quiz.Title = row.Title
		quiz.CreatedAt = row.CreatedAt.Format(time.RFC3339)
		quiz.UpdatedAt = row.UpdatedAt.Format(time.RFC3339)
		quiz.Questions = append(quiz.Questions, views.DBQuestion{
			ID:   strconv.Itoa(int(row.QuestionID)),
			Body: row.QuestionBody,
		})
	}
	return quiz
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
	rows, err := s.Db.Query(r.Context(), `SELECT quizzes.ID, title,
		created_at, updated_at, questions.ID, body FROM quizzes INNER
	JOIN questions ON quizzes.ID = questions.quiz_id WHERE quizzes.ID = $1
	AND owner_id = $2;`, quizID, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var quiz_rows []QuizRow
	for rows.Next() {
		var quiz QuizRow
		err = rows.Scan(&quiz.ID, &quiz.Title, &quiz.CreatedAt,
			&quiz.UpdatedAt, &quiz.QuestionID, &quiz.QuestionBody)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		quiz_rows = append(quiz_rows, quiz)
	}
	quiz := parseRowsToQuiz(quiz_rows)
	page := views.QuizPage(quiz)
	if err := page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
