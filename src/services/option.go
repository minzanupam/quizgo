package services

import (
	"log"
	"net/http"
	"quizgo/src/views"
	"strconv"
)

func (s *Service) optionAddNewApiHandle(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.Atoi(r.PathValue("question_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse question id"))
		return
	}
	optionBody := r.FormValue("option_body")
	if optionBody == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("option body missing"))
		return
	}
	row := s.Db.QueryRow(r.Context(), `INSERT INTO options (body,
		question_id) VALUES ($1, $2) RETURNING ID`, optionBody,
		questionID)
	var option views.DBOption
	if err = row.Scan(&option.ID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	option.Body = optionBody
	compontent := views.OptionCompontent(option)
	if err = compontent.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
