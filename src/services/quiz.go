package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"quizgo/src/views"
	"sort"
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

func (s *Service) quizCreateHandler(w http.ResponseWriter, r *http.Request) {
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
	row := s.DB.QueryRow(r.Context(), `INSERT INTO quizzes (title, owner_id) VALUES ($1, $2) RETURNING ID`, quizTitle, userID)
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
	Status       string
	QuestionID   *int64
	QuestionBody *string
	OptionID     *int64
	OptionBody   *string
}

func parseRowsToQuiz(rows []QuizRow) (views.DBQuiz, error) {
	var quiz views.DBQuiz
	if len(rows) < 1 {
		return views.DBQuiz{}, fmt.Errorf("failed to parse rows with Error: insuffient number of rows")
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].ID != rows[j].ID {
			return rows[i].ID < rows[j].ID
		}
		if rows[i].QuestionID != nil && rows[j].QuestionID != nil && *rows[i].QuestionID != *rows[j].QuestionID {
			return *rows[i].QuestionID < *rows[j].QuestionID
		}
		if rows[i].OptionID != nil && rows[j].OptionID != nil && *rows[i].OptionID != *rows[j].OptionID {
			return *rows[i].OptionID < *rows[j].OptionID
		}
		return false
	})
	row1 := rows[0]
	quiz.ID = strconv.Itoa(int(row1.ID))
	quiz.Title = row1.Title
	quiz.CreatedAt = row1.CreatedAt.Format(time.RFC3339)
	quiz.UpdatedAt = row1.UpdatedAt.Format(time.RFC3339)
	quiz.Status = row1.Status
	if row1.QuestionID == nil {
		return quiz, nil
	}
	pqi := 0 // previous question index
	for i, row := range rows {
		if i > 0 && quiz.Questions[pqi-1].ID == strconv.Itoa(int(*row.QuestionID)) && row.OptionID != nil {
			quiz.Questions[pqi-1].Options = append(quiz.Questions[pqi-1].Options, views.DBOption{
				ID:   strconv.Itoa(int(*row.OptionID)),
				Body: *row.OptionBody,
			})
			continue
		}
		options := make([]views.DBOption, 0)
		if row.OptionID != nil {
			options = append(options, views.DBOption{
				ID:   strconv.Itoa(int(*row.OptionID)),
				Body: *row.OptionBody,
			})
		}
		quiz.Questions = append(quiz.Questions, views.DBQuestion{
			ID:      strconv.Itoa(int(*row.QuestionID)),
			Body:    *row.QuestionBody,
			Options: options,
		})
		pqi += 1
	}
	return quiz, nil
}

func (s *Service) quizPageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(s.Store, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	quizID, err := strconv.Atoi(r.PathValue("quiz_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err := s.DB.Query(r.Context(), `
		SELECT
			quizzes.ID, quizzes.title, quizzes.created_at,
			quizzes.updated_at, quizzes.status, questions.ID,
			questions.body, options.ID, options.body
		FROM
			quizzes
		LEFT JOIN
			questions
		ON
			quizzes.ID = questions.quiz_id
		LEFT JOIN
			options
		ON
			questions.ID = options.question_id
		WHERE
			quizzes.ID = $1
			AND owner_id = $2
	`, quizID, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var quiz_rows []QuizRow
	for rows.Next() {
		var row QuizRow
		err = rows.Scan(&row.ID, &row.Title, &row.CreatedAt,
			&row.UpdatedAt, &row.Status, &row.QuestionID, &row.QuestionBody, &row.OptionID, &row.OptionBody)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		quiz_rows = append(quiz_rows, row)
	}
	quiz, err := parseRowsToQuiz(quiz_rows)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := views.QuizPage(quiz)
	if err := page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) quizPublishHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("failed to parse quiz_id"))
		return
	}

	tx, err := s.DB.Begin(r.Context())
	defer tx.Rollback(context.Background())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	row := s.DB.QueryRow(r.Context(), `SELECT status FROM quizzes WHERE ID = $1 AND owner_id = $2`, quizID, userID)
	var quizStatus string
	if err = row.Scan(&quizStatus); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if quizStatus == "published" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("quiz alread published"))
		return
	} else if quizStatus == "expired" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("quiz has expired"))
		return
	} else if quizStatus != "unpublished" {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("unknow quiz status"))
		return
	}

	_, err = s.DB.Exec(r.Context(), `UPDATE quizzes SET status = 'published' WHERE ID = $1`, quizID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
