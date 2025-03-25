package services

import (
	"log"
	"net/http"
	"quizgo/src/views"
)

func (s *Service) dashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := authorize(s.Store, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login?redirect_url=%2Fdashboard", http.StatusTemporaryRedirect)
	}
	page := views.DashboardPage()
	if err = page.Render(r.Context(), w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
