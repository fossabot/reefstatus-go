package routes

import (
	"encoding/json"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func SetupCommandRoute(r *mux.Router) {
	commandRoute := r.PathPrefix("/command").Subrouter()
	commandRoute.HandleFunc("/feedpasue", handleFeedPasue).Methods("POST")
	commandRoute.HandleFunc("/thunderstorm", handleThunderstorm).Methods("POST")
	commandRoute.HandleFunc("/resetReminder/{Index}", handleResetReminder).Methods("POST")
	commandRoute.HandleFunc("/maintenance/{Index}", handleMaintenance).Methods("POST")
	commandRoute.HandleFunc("/clearlevelalarm/{Id}", handleClearLevelAlarm).Methods("POST")
	commandRoute.HandleFunc("/startwaterchange/{Id}", handleStartWaterChange).Methods("POST")
}

func handleFeedPasue(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handleFeedPasue %s", r.URL.String())
	var body []byte
	r.Body.Read(body)
	var enable bool
	json.Unmarshal(body, &enable)
	data.Controller.FeedPause(enable)
	reply, _ := json.Marshal(true)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleThunderstorm(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handleThunderstorm %s", r.URL.String())
	var body []byte
	r.Body.Read(body)
	var duration int
	json.Unmarshal(body, &duration)
	data.Controller.Thunderstorm(duration)
	reply, _ := json.Marshal(true)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleResetReminder(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handleResetReminder %s", r.URL.String())
	vars := mux.Vars(r)
	index, _ := strconv.Atoi(vars["Index"])
	data.Controller.ResetReminder(index)
	reply, _ := json.Marshal(true)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
func handleMaintenance(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handleMaintenance %s", r.URL.String())
	vars := mux.Vars(r)
	index, _ := strconv.Atoi(vars["Index"])
	var body []byte
	r.Body.Read(body)
	var enable bool
	json.Unmarshal(body, &enable)

	data.Controller.Maintenance(index, enable)
	reply, _ := json.Marshal(true)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
func handleClearLevelAlarm(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handleClearLevelAlarm %s", r.URL.String())
	vars := mux.Vars(r)
	id, _ := vars["Id"]
	data.Controller.ClearLevelAlarm(id)
	reply, _ := json.Marshal(true)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
func handleStartWaterChange(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handleThunderstorm %s", r.URL.String())
	vars := mux.Vars(r)
	id, _ := vars["Id"]
	data.Controller.WaterChange(id)
	reply, _ := json.Marshal(true)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
