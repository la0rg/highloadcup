package store

import (
	"errors"
	"sync"

	"github.com/la0rg/highloadcup/model"
)

var (
	ErrRequiredFields = errors.New("Not all required fields are filled")
	ErrAlreadyExist   = errors.New("Already exist")
	ErrDoesNotExist   = errors.New("Does not exist")
	ErrIDInUpdate     = errors.New("Update should not contain ID in the json object")
)

// Store is an object that keeps all the data (in memory)
// and provides all the aggregation functions
type Store struct {
	mx                 sync.RWMutex
	usersByID          map[int32]*model.User
	visitsByID         map[int32]*model.Visit
	visitsByUserID     map[int32]*VisitIndex
	visitsByLocationID map[int32]*VisitIndex
	locationsByID      map[int32]*model.Location
}

// NewStore constructor
func NewStore() *Store {
	return &Store{
		usersByID:          make(map[int32]*model.User),
		visitsByID:         make(map[int32]*model.Visit),
		visitsByUserID:     make(map[int32]*VisitIndex),
		visitsByLocationID: make(map[int32]*VisitIndex),
		locationsByID:      make(map[int32]*model.Location),
	}
}

// AddUser adds new user to the store
func (s *Store) AddUser(user model.User) error {
	if user.BirthDate == nil || user.Email == nil || user.FirstName == nil ||
		user.LastName == nil || user.Gender == nil || user.ID == nil {
		return ErrRequiredFields
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.usersByID[*(user.ID)]
	if ok {
		return ErrAlreadyExist
	}
	s.usersByID[*(user.ID)] = &user

	// initialize visitsByUserID with empty index (to return [] if user exist and visits were not added)
	vi, ok := s.visitsByUserID[*(user.ID)]
	if ok {
		vi.ApplyToAll(func(visit *model.Visit) {
			visit.User = &user
		})
	} else {
		s.visitsByUserID[*(user.ID)] = NewVisitIndex()
	}
	return nil
}

// GetUserByID find user by id
// returns user and existence flag (map like)
func (s *Store) GetUserByID(id int32) (*model.User, bool) {
	var result model.User
	s.mx.RLock()
	defer s.mx.RUnlock()
	u, ok := s.usersByID[id]
	if ok {
		result = *u // return copy of the object pointed by u
	}
	return &result, ok
}

// UpdateUserByID updates user with id by user
// if user does not exist returns error
func (s *Store) UpdateUserByID(id int32, user model.User) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	u, ok := s.usersByID[id]
	if !ok {
		return ErrDoesNotExist
	}
	if user.ID != nil {
		return ErrIDInUpdate
	}
	if user.BirthDate != nil {
		u.BirthDate = user.BirthDate
	}
	if user.Email != nil {
		u.Email = user.Email
	}
	if user.FirstName != nil {
		u.FirstName = user.FirstName
	}
	if user.LastName != nil {
		u.LastName = user.LastName
	}
	if user.Gender != nil {
		u.Gender = user.Gender
	}
	return nil
}

// AddLocation adds new location to the store
func (s *Store) AddLocation(location model.Location) error {
	if location.City == nil || location.Country == nil || location.Distance == nil ||
		location.ID == nil || location.Place == nil {
		return ErrRequiredFields
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.locationsByID[*(location.ID)]
	if ok {
		return ErrAlreadyExist
	}
	s.locationsByID[*(location.ID)] = &location

	// update connections (if already exist to this entity)
	vi, ok := s.visitsByLocationID[*(location.ID)]
	if ok {
		vi.ApplyToAll(func(visit *model.Visit) {
			visit.Location = &location
		})
	} else {
		// initialize visitsByLocationID with empty index (to return 0 avg)
		s.visitsByLocationID[*(location.ID)] = NewVisitIndex()
	}
	return nil
}

// GetLocationByID find location by id
func (s *Store) GetLocationByID(id int32) (*model.Location, bool) {
	var result model.Location
	s.mx.RLock()
	defer s.mx.RUnlock()
	l, ok := s.locationsByID[id]
	if ok {
		result = *l
	}
	return &result, ok
}

func (s *Store) GetLocationAvg(id int32, fromDate *int64, toDate *int64, fromAge *int64, toAge *int64, gender *string) (float64, bool) {
	var avg float64
	s.mx.RLock()
	defer s.mx.RUnlock()
	vi, ok := s.visitsByLocationID[id]
	if ok {
		visits := vi.GetByAgeAndGender(fromDate, toDate, fromAge, toAge, gender)
		if len(visits) > 0 {
			for i := range visits {
				avg += float64(*(visits[i].Mark))
			}
			avg = avg / float64(len(visits))
		}
		return avg, true
	}
	return 0, false
}

// UpdateLocationByID updates location with id by user
// if location does not exist returns error
func (s *Store) UpdateLocationByID(id int32, location model.Location) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	l, ok := s.locationsByID[id]
	if !ok {
		return ErrDoesNotExist
	}
	if location.ID != nil {
		return ErrIDInUpdate
	}
	if location.City != nil {
		l.City = location.City
	}
	if location.Country != nil {
		l.Country = location.Country
	}
	if location.Distance != nil {
		l.Distance = location.Distance
	}
	if location.Place != nil {
		l.Place = location.Place
	}
	return nil
}

func (s *Store) addVisitToVisitsByLocationID(visit *model.Visit) {
	vi, ok := s.visitsByLocationID[*(visit.LocationID)]
	if !ok {
		vi = NewVisitIndex()
		s.visitsByLocationID[*(visit.LocationID)] = vi
	}
	vi.Add(visit)
}

func (s *Store) addVisitToVisitsByUserID(visit *model.Visit) {
	visitIndex, ok := s.visitsByUserID[*(visit.UserID)]
	if !ok {
		visitIndex = NewVisitIndex()
		s.visitsByUserID[*(visit.UserID)] = visitIndex
	}
	visitIndex.Add(visit)
}

func (s *Store) updateLocationLink(visit *model.Visit) {
	location, ok := s.locationsByID[*(visit.LocationID)]
	if ok {
		visit.Location = location
	}
}

func (s *Store) updateUserLink(visit *model.Visit) {
	user, ok := s.usersByID[*(visit.UserID)]
	if ok {
		visit.User = user
	}
}

// AddVisit adds new visit to the store
func (s *Store) AddVisit(visit model.Visit) error {
	if visit.ID == nil || visit.LocationID == nil ||
		visit.UserID == nil || visit.Mark == nil || visit.VisitedAt == nil {
		return ErrRequiredFields
	}

	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.visitsByID[*(visit.ID)]
	if ok {
		return ErrAlreadyExist
	}
	s.visitsByID[*(visit.ID)] = &visit

	s.addVisitToVisitsByLocationID(&visit)
	s.addVisitToVisitsByUserID(&visit)

	// connect to location
	s.updateLocationLink(&visit)
	// connect to user
	s.updateUserLink(&visit)
	return nil
}

// GetVisitByID find visit by id
func (s *Store) GetVisitByID(id int32) (*model.Visit, bool) {
	var result model.Visit
	s.mx.RLock()
	defer s.mx.RUnlock()
	v, ok := s.visitsByID[id]
	if ok {
		result = *v
	}
	return &result, ok
}

func (s *Store) GetVisitsByUserID(id int32, fromDate *int64, toDate *int64, country *string, toDistance *int32) (*model.UserVisitArray, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	visitIndex, ok := s.visitsByUserID[id]
	if ok {
		visits := visitIndex.GetByCountryAndDistance(fromDate, toDate, country, toDistance)
		return &visits, true
	}
	return nil, false
}

// UpdateVisitByID updates visit with id by visit
// if visit does not exist returns error
func (s *Store) UpdateVisitByID(id int32, visit model.Visit) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	v, ok := s.visitsByID[id]
	if !ok {
		return ErrDoesNotExist
	}
	if visit.ID != nil {
		return ErrIDInUpdate
	}

	if visit.LocationID != nil {
		if *(v.LocationID) != *(visit.LocationID) {
			// remove from the old id position
			vl, ok := s.visitsByLocationID[*(v.LocationID)]
			if ok {
				vl.Remove(v)
			}

			v.LocationID = visit.LocationID
			s.updateLocationLink(v)

			// add again with updated id
			s.addVisitToVisitsByLocationID(v)
		}
	}
	if visit.Mark != nil {
		v.Mark = visit.Mark
	}
	if visit.UserID != nil {
		if *(v.UserID) != *(visit.UserID) {
			// transfer from one VisitIndex to another
			vi, ok := s.visitsByUserID[*(v.UserID)]
			if ok {
				vi.Remove(v)
			}

			v.UserID = visit.UserID
			s.updateUserLink(v)
			s.addVisitToVisitsByUserID(v)
		}
	}
	if visit.VisitedAt != nil {
		if *(v.VisitedAt) != *(visit.VisitedAt) {
			// Delete and insert again visit to the VisitIndex (tree rebalancing)
			vi1, ok1 := s.visitsByUserID[*(v.UserID)]
			vi2, ok2 := s.visitsByLocationID[*(v.LocationID)]
			if ok1 {
				vi1.Remove(v)
			}
			if ok2 {
				vi2.Remove(v)
			}

			v.VisitedAt = visit.VisitedAt

			s.addVisitToVisitsByUserID(v)
			s.addVisitToVisitsByLocationID(v)
		}
	}
	return nil
}
