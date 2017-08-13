package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/la0rg/highloadcup/store"
	"github.com/la0rg/highloadcup/util"
	log "github.com/sirupsen/logrus"
)

var dataStore = store.NewStore()

func main() {
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.NotFound = NotFound

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

func routing(router *httprouter.Router) {
	router.GET("/users/:id", User)
	router.POST("/users/:id", UserUpdate)

	router.GET("/locations/:id", Location)
	router.POST("/locations/:id", LocationUpdate)

	router.GET("/visits/:id", Visit)
	router.POST("/visits/:id", VisitUpdate)
}
