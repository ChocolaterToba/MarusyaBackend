package routing

import (
	"cmkids/interfaces/basic"
	mid "cmkids/interfaces/middleware"

	"github.com/gorilla/mux"
)

func CreateRouter(basicHandler *basic.BasicHandler) *mux.Router {
	r := mux.NewRouter()

	r.Use(mid.PanicMid)

	// r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/api/basic", basicHandler.HandleBasicRequest).Methods("POST", "GET")
	return r
}
