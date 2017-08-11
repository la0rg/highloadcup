package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/la0rg/highloadcup/store"
	"github.com/la0rg/highloadcup/util"
	log "github.com/sirupsen/logrus"
)

var dataStore = store.NewStore()

func main() {
	router := httprouter.New()

	// import static data
	err := util.ImportDataFromZip(dataStore)
	if err != nil {
		log.Fatal(err)
	}

	routing(router)

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	h := &http.Server{Addr: ":80", Handler: router}
	go func() {
		log.Fatal(h.ListenAndServe())
	}()

	<-stop
	log.Info("Shutting down the server...")
	ctx, cn := context.WithTimeout(context.Background(), 2*time.Second)
	defer cn()
	h.Shutdown(ctx)
}

func routing(router *httprouter.Router) {
	router.GET("/users/:id", Users)
}

// Users returns a user by id
func Users(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	strId := params.ByName("id")
	if strId == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	id64, err := strconv.ParseInt(strId, 10, 32)
	id := int32(id64)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	user, ok := dataStore.GetUserByID(id)
	if ok {
		err = writeStructAsJSON(w, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusNotFound)
}

func writeStructAsJSON(w http.ResponseWriter, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(b)))
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
