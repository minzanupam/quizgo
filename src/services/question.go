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

func (s *Service) questionCreateHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(s.Store, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	quizID, err := strconv.Atoi(r.FormValue("quiz_id"))
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
	row := s.DB.QueryRow(r.Context(), `SELECT FROM quizzes WHERE quizzes.ID = $1 and owner_id = $2`, quizID, userID)
	if err = row.Scan(); err != nil {
		log.Println(err)
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	row = s.DB.QueryRow(r.Context(), `INSERT INTO questions (body,
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
	component := views.Question(question)
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
	_, err = s.DB.Exec(context.Background(), `UPDATE questions SET body
	= $1 WHERE quiz_id = $2 and ID = $3`, questionBody, quizID, questionID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) questionEditCompontentHandler(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.Atoi(r.PathValue("question_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rows, err := s.DB.Query(context.Background(), `
		SELECT
			questions.ID, questions.quiz_id, questions.body, options.ID, options.Body
		FROM
			questions
		INNER JOIN
			options
		ON
			options.question_id = questions.ID
		WHERE
			questions.ID = $1
		`, questionID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var question views.DBQuestion
	var quizID int64
	if rows.Next() {
		var option views.DBOption
		if err = rows.Scan(&question.ID, &quizID, &question.Body, &option.ID, &option.Body); err != nil {
			log.Println(err)
		}
		question.Options = []views.DBOption{option}
	}
	for rows.Next() {
		var option views.DBOption
		if err = rows.Scan(nil, nil, nil, &option.ID, &option.Body); err != nil {
			log.Println(err)
		}
		question.Options = append(question.Options, option)
	}
	component := views.QuestionEditComponent(strconv.Itoa(int(quizID)), question)
	if err = component.Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func (s *Service) questionUpdateValuesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}
