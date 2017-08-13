package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Users returns a user by id
func Users(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

// Locations returns a location by id
func Locations(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

// Visits returns a visit by id
func Visits(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
