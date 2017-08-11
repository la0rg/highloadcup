package store

import (
	"sync"

	"github.com/la0rg/highloadcup/model"
)

// Store is an object that keeps all the data (in memory)
// and provides all the aggregation functions
type Store struct {
	mx            sync.Mutex
	usersByID     map[int32]*model.User
	visitsByID    map[int32]*model.Visit
	visitsByUser  map[int32]*model.Visit
	locationsByID map[int32]*model.Location
}

// NewStore constructor
func NewStore() *Store {
	return &Store{
		usersByID:     make(map[int32]*model.User),
		visitsByID:    make(map[int32]*model.Visit),
		visitsByUser:  make(map[int32]*model.Visit),
		locationsByID: make(map[int32]*model.Location),
	}
}

// AddUser adds new user to the store
func (s *Store) AddUser(user model.User) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.usersByID[user.ID] = &user
}

// GetUserByID find user by id
// returns user and existence flag (map like)
func (s *Store) GetUserByID(id int32) (model.User, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	u, ok := s.usersByID[id]
	return *u, ok // return copy of the object pointed by u
}

// AddLocation adds new location to the store
func (s *Store) AddLocation(location model.Location) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.locationsByID[location.ID] = &location
}

// GetLocationByID find location by id
func (s *Store) GetLocationByID(id int32) (model.Location, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	l, ok := s.locationsByID[id]
	return *l, ok
}

// AddVisit adds new visit to the store
func (s *Store) AddVisit(visit model.Visit) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.visitsByID[visit.ID] = &visit
}

// GetVisitByID find visit by id
func (s *Store) GetVisitByID(id int32) (model.Visit, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	v, ok := s.visitsByID[id]
	return *v, ok
}
