package main

import (
	"context"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/frontend/routes"
	"github.com/cjburchell/reefstatus-go/history"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := history.DataInstance.Setup()
	if err != nil {
		log.Fatalf("Unable to Connect to history database", err.Error())
	}

	go UpdateAlerts()
	go UpdateHistory()
	go UpdateWeekHistory()
	go UpdateYearHistory()
	go UpdateController()

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
		Addr:         ":8090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)

	log.Print("shutting down")
	os.Exit(0)
}
