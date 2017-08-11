package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/la0rg/highloadcup/model"
	"github.com/la0rg/highloadcup/util"
	log "github.com/sirupsen/logrus"
)

// usersData is inmemory user storage
var usersData map[int32]model.User

func main() {
	router := httprouter.New()

	// init inmemory store
	usersData = make(map[int32]model.User)

	// import static data
	err := util.ImportDataFromZip(usersData)
	if err != nil {
		log.Fatal(err)
	}

	routing(router)
	log.Fatal(http.ListenAndServe(":80", router))
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
	user, ok := usersData[id]
	if ok {
		err = writeStructAsJSON(w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
