package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/antonlindstrom/pgstore"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Service struct {
	Db    *pgx.Conn
	Store *pgstore.PGStore
}

func (s *Service) rootHandler(w http.ResponseWriter, r *http.Request) {
	page, err := template.ParseFiles("src/views/index.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	row := s.Db.QueryRow(r.Context(), `SELECT 'hello world from database'`)
	var message string
	row.Scan(&message)
	if err = page.Execute(w, map[string]interface{}{"Message": message}); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	store, err := pgstore.NewPGStore(os.Getenv("DATABASE_URL"), []byte(os.Getenv("SESSION_SECRET")))
	s := Service{
		Db:    conn,
		Store: store,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.rootHandler)

	http.ListenAndServe(":4000", mux)
}
