package services

import (
	"fmt"
	"log"
	"net/http"
	"quizgo/src/views"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID       int64
	FullName string
	Email    string
	Password string
}

func (s *Service) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	page := views.LoginPage(nil)
	if err := page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) loginApiHandler(w http.ResponseWriter, r *http.Request) {
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
	row := s.Db.QueryRow(r.Context(), `SELECT ID, email, password FROM users WHERE email = $1`, req_email)
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if err != pgx.ErrNoRows {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		page := views.LoginPage(fmt.Errorf("email not found"))
		if err = page.Render(r.Context(), w); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if user.Password != req_password {
		page := views.LoginPage(fmt.Errorf("incorrect email or password"))
		if err := page.Render(r.Context(), w); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	session, err := s.Store.Get(r, "authsession")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	session.Values["user_id"] = user.ID
	if err = session.Save(r, w); err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
