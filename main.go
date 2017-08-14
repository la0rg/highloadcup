package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/la0rg/highloadcup/store"
	"github.com/la0rg/highloadcup/util"
	log "github.com/sirupsen/logrus"
)

var dataStore = store.NewStore()

func main() {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(NotFound)

	// import static data
	err := util.ImportDataFromZip(dataStore)
	if err != nil {
		log.Fatal(err)
	}

	// set up routes
	routing(router)

	// start http server
	h := &http.Server{Addr: ":80", Handler: router}
	go func() {
		log.Fatal(h.ListenAndServe())
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Info("Shutting down the server...")
	ctx, cn := context.WithTimeout(context.Background(), 2*time.Second)
	defer cn()
	h.Shutdown(ctx)
}

func routing(router *mux.Router) {
	router.HandleFunc("/users/{id}", User).Methods("GET")
	router.HandleFunc("/users/new", UserCreate).Methods("POST")
	router.HandleFunc("/users/{id}", UserUpdate).Methods("POST")

	router.HandleFunc("/locations/{id}", Location).Methods("GET")
	router.HandleFunc("/locations/new", LocationCreate).Methods("POST")
	router.HandleFunc("/locations/{id}", LocationUpdate).Methods("POST")

	router.HandleFunc("/visits/{id}", Visit).Methods("GET")
	router.HandleFunc("/visits/new", VisitCreate).Methods("POST")
	router.HandleFunc("/visits/{id}", VisitUpdate).Methods("POST")
}
