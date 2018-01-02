package routes

import (
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupSettingsRoute(r *mux.Router) {
	settingsRoute := r.PathPrefix("/settings").Subrouter()
	settingsRoute.HandleFunc("/connection", handleSettings).Methods("GET")
	settingsRoute.HandleFunc("/logging", handleSettings).Methods("GET")
	settingsRoute.HandleFunc("/connection", handleSettings).Methods("PUT")
	settingsRoute.HandleFunc("/logging", handleSettings).Methods("PUT")
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handleSettings %s %s", r.Method, r.URL.String())
}
