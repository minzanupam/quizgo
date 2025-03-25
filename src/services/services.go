package services

import (
	"context"
	"log"
	"net/http"
	"os"
	"quizgo/src/views"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Service struct {
	Db    *pgx.Conn
	Store *pgstore.PGStore
}

func (s *Service) rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := authorize(s.Store, r)
	auth := true
	if err != nil {
		log.Println(err)
		auth = false
	}
	page := views.RootPage("hey", auth)
	page.Render(r.Context(), w)
}

func HttpService() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	store, err := pgstore.NewPGStore(os.Getenv("DATABASE_URL"), []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	store.StopCleanup(store.Cleanup(time.Minute * 5))
	s := Service{
		Db:    conn,
		Store: store,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.rootHandler)
	mux.HandleFunc("GET /login", s.loginPageHandler)
	mux.HandleFunc("GET /signup", s.signupPageHandler)
	mux.HandleFunc("GET /profile", s.profilePageHandler)

	mux.HandleFunc("POST /login", s.loginApiHandler)
	mux.HandleFunc("POST /signup", s.signupApiHandler)
	mux.HandleFunc("POST /logout", s.logoutHandler)

	if err = http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
