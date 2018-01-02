package routes

import (
	"encoding/json"
	"github.com/cjburchell/reefstatus-go/data"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupControllerRoute(r *mux.Router) {
	controllerRoute := r.PathPrefix("/controller").Subrouter()
	controllerRoute.HandleFunc("/info", handleInfo).Methods("GET")
	controllerRoute.HandleFunc("/probe", handleProbe).Methods("GET")
	controllerRoute.HandleFunc("/levelsensor", handleLevelSensor).Methods("GET")
	controllerRoute.HandleFunc("/sport", handleSPort).Methods("GET")
	controllerRoute.HandleFunc("/lport", handleLPort).Methods("GET")
	controllerRoute.HandleFunc("/digitalinput", handleDigitalInput).Methods("GET")
	controllerRoute.HandleFunc("/pump", handlePump).Methods("GET")
	controllerRoute.HandleFunc("/programmablelogic", handleProgrammableLogic).Methods("GET")
	controllerRoute.HandleFunc("/dosingpump", handleDosingPump).Methods("GET")
	controllerRoute.HandleFunc("/light", handleLight).Methods("GET")
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetInfo())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleProbe(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetProbes())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleLevelSensor(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetLevelSensors())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleSPort(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetSPorts())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleLPort(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetLPorts())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleDigitalInput(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetDigitalInputs())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handlePump(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetCurrentPumps())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleProgrammableLogic(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetProgrammableLogic())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleDosingPump(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetDosingPumps())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}

func handleLight(w http.ResponseWriter, r *http.Request) {
	reply, _ := json.Marshal(data.Controller.GetLights())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
