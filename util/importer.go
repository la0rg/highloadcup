package util

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/la0rg/highloadcup/model"
	"github.com/la0rg/highloadcup/store"
	log "github.com/sirupsen/logrus"
)

var errUnsupportedFile = errors.New("Imported file is not supported")

func ImportDataFromZip(store *store.Store) error {
	r, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		name := f.Name

		// cut off folder part
		i := strings.LastIndex(name, "/")
		if i != -1 {
			name = name[i:]
		}
		// cut off extension part
		i = strings.LastIndex(name, ".")
		if i != -1 {
			name = name[:i]
		}

		parts := strings.Split(name, "_")
		if len(parts) != 2 {
			return errUnsupportedFile
		}

		log.Infof("Start reading file %s", name)
		// add concurrent processing
		rc, err := f.Open()
		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}
		err = rc.Close()
		if err != nil {
			return err
		}

		switch parts[0] {
		case "users":
			err = importUsers(bytes, store)
		case "locations":
			err = importLocations(bytes, store)
		case "visits":
			err = importVisits(bytes, store)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func importUsers(b []byte, store *store.Store) error {
	var users model.UserArray
	err := json.Unmarshal(b, &users)
	if err != nil {
		return err
	}
	for _, u := range users.Users {
		store.AddUser(u)
	}
	return nil
}

func importLocations(b []byte, store *store.Store) error {
	var locations model.LocationArray
	err := json.Unmarshal(b, &locations)
	if err != nil {
		return err
	}
	for _, l := range locations.Locations {
		store.AddLocation(l)
	}
	return nil
}

func importVisits(b []byte, store *store.Store) error {
	var visits model.VisitArray
	err := json.Unmarshal(b, &visits)
	if err != nil {
		return err
	}
	for _, v := range visits.Visits {
		store.AddVisit(v)
	}
	return nil
}

func ImportCurrentTimestamp() time.Time {
	file, err := os.Open("/tmp/data/options.txt")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	scan := bufio.NewScanner(file)
	if !scan.Scan() {
		log.Fatal("Could not scan the option.txt file.")
	}
	timestampText := scan.Text()
	i64, err := strconv.ParseInt(timestampText, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return time.Unix(i64, 0)
}
