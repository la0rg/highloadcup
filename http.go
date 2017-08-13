package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/la0rg/highloadcup/model"

	"github.com/julienschmidt/httprouter"
)

// User returns a user by id
func User(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseID(params)
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

// UserUpdate update user entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func UserUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseID(params)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user model.User
	errParse := json.Unmarshal(bytes, &user)
	if errParse != nil {
		user = model.User{}
	}
	err = dataStore.UpdateUserByID(id, user)
	// 404 is a higher priority than 400
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Location returns a location by id
func Location(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseID(params)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	location, ok := dataStore.GetLocationByID(id)
	if ok {
		err = writeStructAsJSON(w, &location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusNotFound)
}

// LocationUpdate update location entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func LocationUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseID(params)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var location model.Location
	errParse := json.Unmarshal(bytes, &location)
	if errParse != nil {
		location = model.Location{}
	}
	err = dataStore.UpdateLocationByID(id, location)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Visit returns a visit by id
func Visit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseID(params)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	visit, ok := dataStore.GetVisitByID(id)
	if ok {
		err = writeStructAsJSON(w, &visit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusNotFound)
}

// VisitUpdate update user entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func VisitUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := parseID(params)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var visit model.Visit
	errParse := json.Unmarshal(bytes, &visit)
	if errParse != nil {
		visit = model.Visit{}
	}
	err = dataStore.UpdateVisitByID(id, visit)
	// 404 is a higher priority than 400
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// NotFound custom request handler for non-found requests
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func parseID(params httprouter.Params) (int32, error) {
	errParse := errors.New("Could not parse Id from request")
	strID := params.ByName("id")
	if strID == "" {
		return 0, errParse
	}
	id64, err := strconv.ParseInt(strID, 10, 32)
	id := int32(id64)
	if err != nil {
		return 0, errParse
	}
	return id, nil
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
