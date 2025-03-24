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
	row := s.Db.QueryRow(r.Context(), `SELECT 'hello world from database'`)
	var message string
	row.Scan(&message)
	page := views.RootPage(message)
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
	mux.HandleFunc("GET /login", s.loginPageHandler)
	mux.HandleFunc("GET /signup", s.signupPageHandler)
	mux.HandleFunc("GET /", s.rootHandler)

	if err = http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
