package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"quizgo/src/views"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func (s *Service) questionBlockHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Service) questionApiAddHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(s.Store, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	quizID, err := strconv.Atoi(r.PathValue("quiz_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	questionBody := r.FormValue("question_body")
	if questionBody == "" {
		log.Println(fmt.Errorf("empty question body"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row := s.Db.QueryRow(r.Context(), `SELECT FROM quizzes WHERE quizzes.ID = $1 and owner_id = $2`, quizID, userID)
	if err = row.Scan(); err != nil {
		log.Println(err)
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	row = s.Db.QueryRow(r.Context(), `INSERT INTO questions (body,
		quiz_id) VALUES ($1, $2) RETURNING ID`, questionBody, quizID)
	var questionID int64
	if err = row.Scan(&questionID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	question := views.DBQuestion{
		ID:   strconv.Itoa(int(questionID)),
		Body: questionBody,
	}
	component := views.Question(strconv.Itoa(quizID), question)
	component.Render(r.Context(), w)
}

func (s *Service) questionUpdateNameHandle(w http.ResponseWriter, r *http.Request) {
	questionBody := r.FormValue("question_body")
	quizID, err := strconv.Atoi(r.PathValue("quiz_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	questionID, err := strconv.Atoi(r.PathValue("question_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if questionBody == "" {
		log.Println("empty question body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = s.Db.Exec(context.Background(), `UPDATE questions SET body
	= $1 WHERE quiz_id = $2 and ID = $3`, questionBody, quizID, questionID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
