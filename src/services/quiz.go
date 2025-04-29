package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"quizgo/queries"
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
	quizID, err := s.Queries.CreateQuiz(r.Context(), queries.CreateQuizParams{
		Title:   quizTitle,
		OwnerID: int32(userID),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/dashboard/quiz/%d", int(quizID)), http.StatusFound)
}

func parseRowsToQuiz(rows []queries.GetQuizRow) (views.DBQuiz, error) {
	var quiz views.DBQuiz
	if len(rows) < 1 {
		return views.DBQuiz{}, fmt.Errorf("failed to parse rows with Error: insuffient number of rows")
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].ID != rows[j].ID {
			return rows[i].ID < rows[j].ID
		}
		if rows[i].ID_2.Valid && rows[j].ID_2.Valid && rows[i].ID_2.Int32 != rows[j].ID_2.Int32 {
			return rows[i].ID_2.Int32 < rows[j].ID_2.Int32
		}
		if rows[i].ID_3.Valid && rows[j].ID_3.Valid && rows[i].ID_3.Int32 != rows[j].ID_3.Int32 {
			return rows[i].ID_3.Int32 < rows[j].ID_3.Int32
		}
		return false
	})
	row1 := rows[0]
	quiz.ID = strconv.Itoa(int(row1.ID))
	quiz.Title = row1.Title
	quiz.CreatedAt = row1.CreatedAt.Time.Format(time.RFC3339)
	quiz.UpdatedAt = row1.UpdatedAt.Time.Format(time.RFC3339)
	if err := row1.Status.Scan(quiz.Status); err != nil {
		return views.DBQuiz{}, err
	}
	if !row1.ID_2.Valid {
		log.Println(fmt.Errorf("failed to parse question id"))
		return quiz, nil
	}
	pqi := 0 // previous question index
	for i, row := range rows {
		if i > 0 && quiz.Questions[pqi-1].ID == strconv.Itoa(int(row.ID_2.Int32)) && row.ID_3.Valid {
			quiz.Questions[pqi-1].Options = append(quiz.Questions[pqi-1].Options, views.DBOption{
				ID:   strconv.Itoa(int(row.ID_3.Int32)),
				Body: row.Body_2.String,
			})
			continue
		}
		options := make([]views.DBOption, 0)
		if row.ID_3.Valid {
			options = append(options, views.DBOption{
				ID:   strconv.Itoa(int(row.ID_3.Int32)),
				Body: row.Body_2.String,
			})
		}
		quiz.Questions = append(quiz.Questions, views.DBQuestion{
			ID:      strconv.Itoa(int(row.ID_2.Int32)),
			Body:    row.Body.String,
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

	rows, err := s.Queries.GetQuiz(r.Context(), queries.GetQuizParams{
		ID:      int32(quizID),
		OwnerID: int32(userID),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	quiz, err := parseRowsToQuiz(rows)
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
	row, err := s.Queries.GetQuizStatus(r.Context(), queries.GetQuizStatusParams{
		ID:      int32(quizID),
		OwnerID: int32(userID),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var quizStatus string
	if err = row.Scan(quizStatus); err != nil {
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

	if err = s.Queries.UpdateQuizStatusPublish(r.Context(), int32(quizID)); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
