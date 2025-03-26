package services

import (
	"log"
	"net/http"
	"quizgo/src/views"
	"time"
)

func (s *Service) dashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(s.Store, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login?redirect_url=%2Fdashboard", http.StatusTemporaryRedirect)
	}
	var user views.DBUser
	row := s.Db.QueryRow(r.Context(), `SELECT ID, fullname, email FROM users WHERE ID = $1`, userID)
	if err = row.Scan(&user.ID, &user.FullName, &user.Email); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rows, err := s.Db.Query(r.Context(), `SELECT quizzes.ID, title, 
		created_at, updated_at FROM quizzes WHERE owner_id = $1`, userID)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var quizzes []views.DBQuiz
	for rows.Next() {
		var quiz views.DBQuiz
		var createdAt time.Time
		var updatedAt time.Time
		if err = rows.Scan(&quiz.ID, &quiz.Title, &createdAt, &updatedAt); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		quiz.Owner = user
		quiz.CreatedAt = createdAt.Format(time.RFC3339)
		quiz.UpdatedAt = updatedAt.Format(time.RFC3339)
		quizzes = append(quizzes, quiz)
	}
	page := views.DashboardPage(quizzes)
	if err = page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
