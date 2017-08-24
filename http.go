package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/la0rg/highloadcup/store"
	"github.com/la0rg/highloadcup/util"

	"github.com/gorilla/mux"
	"github.com/la0rg/highloadcup/model"
)

// User returns a user by id
func User(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	user, ok := dataStore.GetUserByID(id)
	if ok {
		err = writeStructAsJSON(w, user)
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
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
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
		if err == store.ErrDoesNotExist {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// UserCreate create user entity
// success - 200 with body {}
// already exist - 400
// incorrect request - 400
func UserCreate(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user model.User
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = dataStore.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// Location returns a location by id
func Location(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	location, ok := dataStore.GetLocationByID(id)
	if ok {
		err = writeStructAsJSON(w, location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusNotFound)
}

// Location returns a location by id
func LocationAvg(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var fromDate, toDate, fromAge, toAge *int64
	var gender *string
	keys, ok := r.URL.Query()["fromDate"]
	if ok && len(keys) >= 1 {
		i64, err := strconv.ParseInt(keys[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fromDate = &i64
	}
	keys, ok = r.URL.Query()["toDate"]
	if ok && len(keys) >= 1 {
		i64, err := strconv.ParseInt(keys[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		toDate = &i64
	}
	keys, ok = r.URL.Query()["fromAge"]
	if ok && len(keys) >= 1 {
		i, err := strconv.Atoi(keys[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		date := now.AddDate(-i, 0, 0).Unix()
		fromAge = &date
	}
	keys, ok = r.URL.Query()["toAge"]
	if ok && len(keys) >= 1 {
		i, err := strconv.Atoi(keys[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		date := now.AddDate(-i, 0, 0).Unix()
		toAge = &date
	}
	keys, ok = r.URL.Query()["gender"]
	if ok && len(keys) >= 1 {
		gender = &(keys[0])
		if !util.IsGender(*gender) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	avg, ok := dataStore.GetLocationAvg(id, fromDate, toDate, fromAge, toAge, gender)
	if ok {
		//avg = util.RoundPlus(avg, 5)
		err = writeStructAsJSON(w, model.Avg{Value: model.FloatPrec5(avg)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// LocationUpdate update location entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func LocationUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
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
		if err == store.ErrDoesNotExist {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// LocationCreate create location entity
// success - 200 with body {}
// already exist - 400
// incorrect request - 400
func LocationCreate(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var location model.Location
	err = json.Unmarshal(bytes, &location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = dataStore.AddLocation(location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// Visit returns a visit by id
func Visit(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	visit, ok := dataStore.GetVisitByID(id)
	if ok {
		err = writeStructAsJSON(w, visit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusNotFound)
}

func VisitsByUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	// // check if user does not exist
	// _, ok := dataStore.GetUserByID(id)
	// if !ok {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	var fromDate, toDate *int64
	var country *string
	var toDistance *int32
	keys, ok := r.URL.Query()["fromDate"]
	if ok && len(keys) >= 1 {
		i64, err := strconv.ParseInt(keys[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fromDate = &i64
	}
	keys, ok = r.URL.Query()["toDate"]
	if ok && len(keys) >= 1 {
		i64, err := strconv.ParseInt(keys[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		toDate = &i64
	}
	keys, ok = r.URL.Query()["country"]
	if ok && len(keys) >= 1 {
		country = &(keys[0])
		if !util.OnlyLetters(*country) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	keys, ok = r.URL.Query()["toDistance"]
	if ok && len(keys) >= 1 {
		i64, err := strconv.ParseInt(keys[0], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		i32 := int32(i64)
		toDistance = &i32
	}

	visits, ok := dataStore.GetVisitsByUserID(id, fromDate, toDate, country, toDistance)
	if ok {
		err = writeStructAsJSON(w, visits)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// VisitUpdate update user entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func VisitUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
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
		if err == store.ErrDoesNotExist {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// VisitCreate create location entity
// success - 200 with body {}
// already exist - 400
// incorrect request - 400
func VisitCreate(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var visit model.Visit
	err = json.Unmarshal(bytes, &visit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = dataStore.AddVisit(visit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// NotFound custom request handler for non-found requests
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func parseID(r *http.Request) (int32, error) {
	errParse := errors.New("Could not parse Id from request")
	vars := mux.Vars(r)
	strID, ok := vars["id"]
	if !ok {
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
