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
	_, err := authenticate(s.Store, r)
	auth := true
	if err != nil {
		log.Println(err)
		auth = false
	}
	page := views.RootPage("hey", auth)
	page.Render(r.Context(), w)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
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

	fs := http.FileServer(http.Dir("public"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", s.rootHandler)
	mux.HandleFunc("GET /login", s.loginPageHandler)
	mux.HandleFunc("GET /signup", s.signupPageHandler)
	mux.HandleFunc("GET /profile", s.profilePageHandler)
	mux.HandleFunc("GET /dashboard", s.dashboardPageHandler)
	mux.HandleFunc("GET /dashboard/quiz", s.quizParentPageHandler)
	mux.HandleFunc("GET /dashboard/quiz/{id}", s.quizPageHandler)
	mux.HandleFunc("GET /dashboard/quiz/question/{question_id}/edit", s.questionEditPage)

	mux.HandleFunc("POST /login", s.loginApiHandler)
	mux.HandleFunc("POST /signup", s.signupApiHandler)
	mux.HandleFunc("POST /logout", s.logoutHandler)
	mux.HandleFunc("POST /quiz", s.quizApiHandler)
	mux.HandleFunc("POST /quiz/{quiz_id}/question", s.questionApiAddHandler)
	mux.HandleFunc("POST /quiz/{quiz_id}/{question_id}/option", s.optionAddNewApiHandle)

	mux.HandleFunc("PUT /quiz/{quiz_id}/question/{question_id}", s.questionUpdateNameHandle)
	mux.HandleFunc("PUT /dashboard/quiz/question/{question_id}/edit", s.questionUpdateValuesHandler)

	if err = http.ListenAndServe(":4000", LoggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
