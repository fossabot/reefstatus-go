package routes

import (
	"encoding/json"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/history"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupDataRoute(r *mux.Router) {
	dataRoute := r.PathPrefix("/data").Subrouter()
	dataRoute.HandleFunc("/log/{Id}", handleDayData).Methods("GET")
	dataRoute.HandleFunc("/logYear/{Id}", handleYearData).Methods("GET")
	dataRoute.HandleFunc("/logWeek/{Id}", handleWeekData).Methods("GET")
}

func handleDayData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["Id"]

	result, err := history.DataInstance.GetDayDataPoints(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened! " + err.Error()))
		return
	}

	reply, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleWeekData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["Id"]

	result, err := history.DataInstance.GetWeekDataPoints(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened! " + err.Error()))
		return
	}

	reply, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleYearData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["Id"]

	result, err := history.DataInstance.GetYearDataPoints(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened! " + err.Error()))
		return
	}

	reply, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
