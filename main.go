package main

import (
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/frontend/routes"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	go Update()

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/web/public/default.html")
	})

	routes.SetupControllerRoute(r)
	routes.SetupCommandRoute(r)
	routes.SetupDataRoute(r)
	routes.SetupSettingsRoute(r)

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("frontend/web/public"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8082",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
