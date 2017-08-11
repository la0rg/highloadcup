package util

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/la0rg/highloadcup/model"
	log "github.com/sirupsen/logrus"
)

var errUnsupportedFile = errors.New("Imported file is not supported")

func ImportDataFromZip(usersData map[int32]model.User) error {
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
		log.Infof("File content is %v", string(bytes))

		switch parts[0] {
		case "users":
			err = importUsers(bytes, usersData)
			if err != nil {
				return err
			}
		case "locations":
			//parse locations
		case "visits":
			// parse locations
		}
	}
	return nil
}

func importUsers(b []byte, usersData map[int32]model.User) error {
	var users model.UserArray
	err := json.Unmarshal(b, &users)
	if err != nil {
		return err
	}
	for _, u := range users.Users {
		usersData[u.ID] = u
	}
	return nil
}
