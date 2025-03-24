package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Service struct {
	Db *pgx.Conn
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
	s := Service {
		Db: conn,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.rootHandler)

	http.ListenAndServe(":4000", mux)
}
