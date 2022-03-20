package routing

import (
	"cmkids/interfaces/basic"
	mid "cmkids/interfaces/middleware"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func CreateRouter(basicHandler *basic.BasicHandler, csrfOn bool, httpsOn bool, ) *mux.Router {
	r := mux.NewRouter()

	r.Use(mid.PanicMid)

	if csrfOn {
		csrfMid := csrf.Protect(
			[]byte(os.Getenv("CSRF_KEY")),
			csrf.Path("/"),
			csrf.Secure(httpsOn),
		)
		r.Use(csrfMid)
		r.Use(mid.CSRFSettingMid)
	}

	// r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/api/basic", basicHandler.HandleBasicRequest).Methods("POST", "GET")

	if csrfOn {
		r.HandleFunc("/api/csrf", func(w http.ResponseWriter, r *http.Request) { // Is used only for getting csrf key
			w.WriteHeader(http.StatusCreated)
		}).Methods("GET")
	}

	return r
}
