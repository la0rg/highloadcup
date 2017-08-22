package store

import (
	"errors"
	"sync"
	"time"

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
	visitsByLocationID map[int32][]*model.Visit
	locationsByID      map[int32]*model.Location
}

// NewStore constructor
func NewStore() *Store {
	return &Store{
		usersByID:          make(map[int32]*model.User),
		visitsByID:         make(map[int32]*model.Visit),
		visitsByUserID:     make(map[int32]*VisitIndex),
		visitsByLocationID: make(map[int32][]*model.Visit),
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
	visits, ok := s.visitsByLocationID[*(location.ID)]
	if ok {
		for i := range visits {
			if visits[i].LocationID == location.ID {
				visits[i].Location = &location
			}
		}
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

	visits, ok := s.visitsByLocationID[*(visit.LocationID)]
	if !ok {
		visits = make([]*model.Visit, 0)
	}
	visits = append(visits, &visit)
	s.visitsByLocationID[*(visit.LocationID)] = visits

	visitIndex, ok := s.visitsByUserID[*(visit.UserID)]
	if !ok {
		visitIndex = NewVisitIndex()
		s.visitsByUserID[*(visit.UserID)] = visitIndex
	}
	visitIndex.Add(&visit)

	// connect to location
	// TODO: check if location might be create after the visit that reference it
	location, ok := s.locationsByID[*(visit.LocationID)]
	if ok {
		visit.Location = location
	}
	// connect to user
	// TODO: -||-
	user, ok := s.usersByID[*(visit.UserID)]
	if ok {
		visit.User = user
	}
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

func (s *Store) GetVisitsByUserID(id int32, fromDate *time.Time, toDate *time.Time, country *string) (*model.VisitArray, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	visitIndex, ok := s.visitsByUserID[id]
	if ok {
		visits := visitIndex.Get(fromDate, toDate, country)
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
		v.LocationID = visit.LocationID
	}
	if visit.Mark != nil {
		v.Mark = visit.Mark
	}
	if visit.UserID != nil {
		v.UserID = visit.UserID
	}
	if visit.VisitedAt != nil {
		v.VisitedAt = visit.VisitedAt
	}
	return nil
}
