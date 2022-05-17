package routing

import (
	"cmkids/interfaces/marusia"
	mid "cmkids/interfaces/middleware"
	"cmkids/interfaces/quiz"

	"github.com/gorilla/mux"
)

func CreateRouter(marusiaHandler *marusia.MarusiaHandler, quizHandler *quiz.QuizHandler) *mux.Router {
	r := mux.NewRouter()

	r.Use(mid.PanicMid)

	// r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/api/marusia", marusiaHandler.HandleMarusiaRequest).Methods("POST", "GET")
	r.HandleFunc("/api/test/add", quizHandler.HandleAddQuiz).Methods("PUT")
	return r
}
