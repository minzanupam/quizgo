package services

import (
	"log"
	"net/http"
	"quizgo/src/views"
)

func (s *Service) signupPageHandler(w http.ResponseWriter, r *http.Request) {
	page := views.SignupPage()
	if err := page.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) signupApiHandler(w http.ResponseWriter, r *http.Request) {
	req_fullname := r.FormValue("fullname")
	if req_fullname == "" {
		w.Write([]byte("full name missing"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req_email := r.FormValue("email")
	if req_email == "" {
		w.Write([]byte("email missing"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req_password := r.FormValue("password")
	if req_password == "" {
		w.Write([]byte("password missing"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row := s.DB.QueryRow(r.Context(), `INSERT INTO users (fullname, email,
		password) VALUES ($1, $2, $3) RETURNING ID`, req_fullname,
		req_email, req_password)
	var userID int64
	if err := row.Scan(&userID); err != nil {
		log.Println(err)
	}
	session, err := s.Store.Get(r, "authsession")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	session.Values["user_id"] = userID
	if err = session.Save(r, w); err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
