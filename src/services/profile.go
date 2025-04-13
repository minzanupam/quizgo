package services

import (
	"fmt"
	"log"
	"net/http"
	"quizgo/src/views"
	"strconv"

	"github.com/antonlindstrom/pgstore"
)

func convertToDbUser(user User) views.DBUser {
	var db_user views.DBUser
	db_user.ID = strconv.Itoa(int(user.ID))
	db_user.FullName = user.FullName
	db_user.Email = user.Email
	return db_user
}

func authenticate(store *pgstore.PGStore, r *http.Request) (int64, error) {
	session, err := store.Get(r, "authsession")
	if err != nil {
		return 0, err
	}
	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		return 0, fmt.Errorf("failed to parse user id")
	}
	return userID, nil
}

func (s *Service) profilePageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(s.Store, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	if userID == 0 {
		log.Println("invalid user id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	row := s.DB.QueryRow(r.Context(), `SELECT ID, fullname, email FROM users WHERE ID = $1`, userID)
	var user views.DBUser
	if err := row.Scan(&user.ID, &user.FullName, &user.Email); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := views.ProfilePage(user)
	if err := page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
