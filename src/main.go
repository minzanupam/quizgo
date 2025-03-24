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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	page, err := template.ParseFiles("src/views/index.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = page.Execute(w, nil); err != nil {
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
	mux := http.NewServeMux()
	log.Println(conn)

	mux.HandleFunc("GET /", rootHandler)

	http.ListenAndServe(":4000", mux)
}
