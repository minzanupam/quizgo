package services

import (
	"fmt"
	"log"
	"net/http"
	"quizgo/queries"
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
	_, err = s.DB.Exec(r.Context(), `UPDATE questions SET body
	= $1 WHERE quiz_id = $2 and ID = $3`, questionBody, quizID, questionID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseQuestion(rows []queries.GetQuestionRow) (views.DBQuestion, error) {
	if len(rows) == 0 {
		return views.DBQuestion{}, fmt.Errorf("failed to parse question: no rows")
	}
	var question views.DBQuestion
	question.ID = strconv.Itoa(int(rows[0].ID))
	question.Body = rows[0].Body
	var options []views.DBOption
	for _, row := range rows {
		var option views.DBOption
		option.ID = strconv.Itoa(int(row.ID_2))
		option.Body = row.Body_2
		options = append(options, option)
	}
	return question, nil
}

func (s *Service) questionEditCompontentHandler(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.Atoi(r.PathValue("question_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rows, err := s.Queries.GetQuestion(r.Context(), int32(questionID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(rows) == 0 {
		log.Println(fmt.Errorf("failed to parse rows: no rows"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	question, err := parseQuestion(rows)
	quizID := rows[0].QuizID
	component := views.QuestionEditComponent(strconv.Itoa(int(quizID)), question)
	if err = component.Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}

func (s *Service) questionUpdateValuesHandler(w http.ResponseWriter, r *http.Request) {
	type Question struct {
		ID   int64
		Body string
	}
	var req_question Question
	var err error
	req_question.ID, err = strconv.ParseInt(r.FormValue("question_id"), 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req_question.Body = r.FormValue("question_body")
	if req_question.Body == "" {
		log.Println("question body is empty")
		w.WriteHeader(http.StatusNotModified)
		return
	}
	if _, err = s.DB.Exec(r.Context(), `
		UPDATE
			questions
		SET
			question.body = $1
		WHERE
			question.ID = $2
		`, req_question.Body, req_question.ID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
