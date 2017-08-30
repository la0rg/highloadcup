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
	mxUsersByID          sync.RWMutex
	mxVisitsByID         sync.RWMutex
	mxVisitsByUserID     sync.RWMutex
	mxVisitsByLocationID sync.RWMutex
	mxLocationsByID      sync.RWMutex
	visitsByUserID       [1100000]*VisitIndex
	usersByID            [1100000]*model.User
	visitsByID           [11000000]*model.Visit
	visitsByLocationID   [1000000]*VisitIndex
	locationsByID        [1000000]*model.Location
}

// NewStore constructor
func NewStore() *Store {
	return &Store{}
}

// AddUser adds new user to the store
func (s *Store) AddUser(user model.User) error {
	if !user.BirthDate.Defined || !user.Email.Defined || !user.FirstName.Defined ||
		!user.LastName.Defined || !user.Gender.Defined || !user.ID.Defined {
		return ErrRequiredFields
	}
	s.mxUsersByID.Lock()
	if user.ID.V >= int32(len(s.usersByID)) || user.ID.V < 0 {
		s.mxUsersByID.Unlock()
		return ErrDoesNotExist
	}
	if s.usersByID[user.ID.V] != nil {
		s.mxUsersByID.Unlock()
		return ErrAlreadyExist
	}
	s.usersByID[user.ID.V] = &user
	s.mxUsersByID.Unlock()

	s.mxVisitsByUserID.Lock()
	// initialize visitsByUserID with empty index (to return [] if user exist and visits were not added)
	vi := s.visitsByUserID[user.ID.V]
	if vi != nil {
		vi.ApplyToAll(func(visit *model.Visit) {
			visit.User = &user
		})
	} else {
		s.visitsByUserID[user.ID.V] = NewVisitIndex()
	}
	s.mxVisitsByUserID.Unlock()
	return nil
}

// GetUserByID find user by id
// returns user and existence flag (map like)
func (s *Store) GetUserByID(id int32) (*model.User, bool) {
	var result model.User
	s.mxUsersByID.RLock()
	defer s.mxUsersByID.RUnlock()
	if id >= int32(len(s.usersByID)) || id < 0 {
		return nil, false
	}
	u := s.usersByID[id]
	if u != nil {
		result = *u // return copy of the object pointed by u
	}
	return &result, u != nil
}

// UpdateUserByID updates user with id by user
// if user does not exist returns error
func (s *Store) UpdateUserByID(id int32, user model.User) error {
	s.mxUsersByID.Lock()
	defer s.mxUsersByID.Unlock()
	if id >= int32(len(s.usersByID)) || id < 0 {
		return ErrDoesNotExist
	}
	u := s.usersByID[id]
	if u == nil {
		return ErrDoesNotExist
	}
	if user.ID.Defined {
		return ErrIDInUpdate
	}
	if user.BirthDate.Defined {
		u.BirthDate = user.BirthDate
	}
	if user.Email.Defined {
		u.Email = user.Email
	}
	if user.FirstName.Defined {
		u.FirstName = user.FirstName
	}
	if user.LastName.Defined {
		u.LastName = user.LastName
	}
	if user.Gender.Defined {
		u.Gender = user.Gender
	}
	return nil
}

// AddLocation adds new location to the store
func (s *Store) AddLocation(location model.Location) error {
	if !location.City.Defined || !location.Country.Defined || !location.Distance.Defined ||
		!location.ID.Defined || !location.Place.Defined {
		return ErrRequiredFields
	}
	s.mxLocationsByID.Lock()
	if location.ID.V >= int32(len(s.locationsByID)) || location.ID.V < 0 {
		s.mxLocationsByID.Unlock()
		return ErrDoesNotExist
	}
	if s.locationsByID[location.ID.V] != nil {
		s.mxLocationsByID.Unlock()
		return ErrAlreadyExist
	}
	s.locationsByID[location.ID.V] = &location
	s.mxLocationsByID.Unlock()

	s.mxVisitsByLocationID.Lock()
	// update connections (if already exist to this entity)
	vi := s.visitsByLocationID[location.ID.V]
	if vi != nil {
		vi.ApplyToAll(func(visit *model.Visit) {
			visit.Location = &location
		})
	} else {
		// initialize visitsByLocationID with empty index (to return 0 avg)
		s.visitsByLocationID[location.ID.V] = NewVisitIndex()
	}
	s.mxVisitsByLocationID.Unlock()
	return nil
}

// GetLocationByID find location by id
func (s *Store) GetLocationByID(id int32) (*model.Location, bool) {
	var result model.Location
	s.mxLocationsByID.RLock()
	defer s.mxLocationsByID.RUnlock()
	if id >= int32(len(s.locationsByID)) || id < 0 {
		return nil, false
	}
	l := s.locationsByID[id]
	if l != nil {
		result = *l
	}
	return &result, l != nil
}

func (s *Store) GetLocationAvg(id int32, fromDate *int64, toDate *int64, fromAge *int64, toAge *int64, gender *string) (float64, bool) {
	var avg float64
	s.mxVisitsByLocationID.RLock()
	defer s.mxVisitsByLocationID.RUnlock()
	if id >= int32(len(s.visitsByLocationID)) || id < 0 {
		return 0, false
	}
	vi := s.visitsByLocationID[id]
	if vi != nil {
		visits := vi.GetByAgeAndGender(fromDate, toDate, fromAge, toAge, gender)
		if len(visits) > 0 {
			for i := range visits {
				avg += float64(visits[i].Mark.V)
			}
			avg = avg / float64(len(visits))
		}
		return avg + 1e-10, true
	}
	return 0, false
}

// UpdateLocationByID updates location with id by user
// if location does not exist returns error
func (s *Store) UpdateLocationByID(id int32, location model.Location) error {
	s.mxLocationsByID.Lock()
	defer s.mxLocationsByID.Unlock()
	if id >= int32(len(s.locationsByID)) || id < 0 {
		return ErrDoesNotExist
	}
	l := s.locationsByID[id]
	if l == nil {
		return ErrDoesNotExist
	}
	if location.ID.Defined {
		return ErrIDInUpdate
	}
	if location.City.Defined {
		l.City = location.City
	}
	if location.Country.Defined {
		l.Country = location.Country
	}
	if location.Distance.Defined {
		l.Distance = location.Distance
	}
	if location.Place.Defined {
		l.Place = location.Place
	}
	return nil
}

func (s *Store) addVisitToVisitsByLocationID(visit *model.Visit) {
	s.mxVisitsByLocationID.Lock()
	defer s.mxVisitsByLocationID.Unlock()
	if visit.LocationID.V >= int32(len(s.visitsByLocationID)) || visit.LocationID.V < 0 {
		return
	}
	vi := s.visitsByLocationID[visit.LocationID.V]
	if vi == nil {
		vi = NewVisitIndex()
		s.visitsByLocationID[visit.LocationID.V] = vi
	}
	vi.Add(visit)
}

func (s *Store) addVisitToVisitsByUserID(visit *model.Visit) {
	s.mxVisitsByUserID.Lock()
	defer s.mxVisitsByUserID.Unlock()
	if visit.UserID.V >= int32(len(s.visitsByUserID)) || visit.UserID.V < 0 {
		return
	}
	visitIndex := s.visitsByUserID[visit.UserID.V]
	if visitIndex == nil {
		visitIndex = NewVisitIndex()
		s.visitsByUserID[visit.UserID.V] = visitIndex
	}
	visitIndex.Add(visit)
}

func (s *Store) updateLocationLink(visit *model.Visit) {
	s.mxLocationsByID.Lock()
	defer s.mxLocationsByID.Unlock()
	if visit.LocationID.V >= int32(len(s.locationsByID)) || visit.LocationID.V < 0 {
		return
	}
	location := s.locationsByID[visit.LocationID.V]
	if location != nil {
		visit.Location = location
	}
}

func (s *Store) updateUserLink(visit *model.Visit) {
	s.mxUsersByID.Lock()
	defer s.mxUsersByID.Unlock()
	if visit.UserID.V >= int32(len(s.usersByID)) || visit.UserID.V < 0 {
		return
	}
	user := s.usersByID[visit.UserID.V]
	if user != nil {
		visit.User = user
	}
}

// AddVisit adds new visit to the store
func (s *Store) AddVisit(visit model.Visit) error {
	if !visit.ID.Defined || !visit.LocationID.Defined ||
		!visit.UserID.Defined || !visit.Mark.Defined || !visit.VisitedAt.Defined {
		return ErrRequiredFields
	}

	s.mxVisitsByID.Lock()
	defer s.mxVisitsByID.Unlock()
	if visit.ID.V >= int32(len(s.visitsByID)) || visit.ID.V < 0 {
		return ErrDoesNotExist
	}

	if s.visitsByID[visit.ID.V] != nil {
		return ErrAlreadyExist
	}
	s.visitsByID[visit.ID.V] = &visit

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
	s.mxVisitsByID.RLock()
	defer s.mxVisitsByID.RUnlock()
	if id >= int32(len(s.visitsByID)) || id < 0 {
		return nil, false
	}
	v := s.visitsByID[id]
	if v != nil {
		result = *v
	}
	return &result, v != nil
}

func (s *Store) GetVisitsByUserID(id int32, fromDate *int64, toDate *int64, country *string, toDistance *int32) (*model.UserVisitArray, bool) {
	s.mxVisitsByUserID.RLock()
	defer s.mxVisitsByUserID.RUnlock()
	if id >= int32(len(s.visitsByUserID)) || id < 0 {
		return nil, false
	}
	visitIndex := s.visitsByUserID[id]
	if visitIndex != nil {
		visits := visitIndex.GetByCountryAndDistance(fromDate, toDate, country, toDistance)
		return &visits, true
	}
	return nil, false
}

// UpdateVisitByID updates visit with id by visit
// if visit does not exist returns error
func (s *Store) UpdateVisitByID(id int32, visit model.Visit) error {
	s.mxVisitsByID.Lock()
	defer s.mxVisitsByID.Unlock()
	if id >= int32(len(s.visitsByID)) || id < 0 {
		return ErrDoesNotExist
	}
	v := s.visitsByID[id]
	if v == nil {
		return ErrDoesNotExist
	}
	if visit.ID.Defined {
		return ErrIDInUpdate
	}

	if visit.LocationID.Defined {
		if v.LocationID.V != visit.LocationID.V {
			// remove from the old id position
			s.mxVisitsByLocationID.Lock()
			vl := s.visitsByLocationID[v.LocationID.V]
			if vl != nil {
				vl.Remove(v)
			}
			s.mxVisitsByLocationID.Unlock()

			v.LocationID = visit.LocationID
			s.updateLocationLink(v)

			// add again with updated id
			s.addVisitToVisitsByLocationID(v)
		}
	}
	if visit.Mark.Defined {
		v.Mark = visit.Mark
	}
	if visit.UserID.Defined {
		if v.UserID.V != visit.UserID.V {
			s.mxVisitsByUserID.Lock()
			// transfer from one VisitIndex to another
			vi := s.visitsByUserID[v.UserID.V]
			if vi != nil {
				vi.Remove(v)
			}
			s.mxVisitsByUserID.Unlock()

			v.UserID = visit.UserID
			s.updateUserLink(v)
			s.addVisitToVisitsByUserID(v)
		}
	}
	if visit.VisitedAt.Defined {
		if v.VisitedAt.V != visit.VisitedAt.V {
			// Delete and insert again visit to the VisitIndex (tree rebalancing)
			s.mxVisitsByLocationID.Lock()
			s.mxVisitsByUserID.Lock()
			vi1 := s.visitsByUserID[v.UserID.V]
			vi2 := s.visitsByLocationID[v.LocationID.V]
			if vi1 != nil {
				vi1.Remove(v)
			}
			if vi2 != nil {
				vi2.Remove(v)
			}
			s.mxVisitsByLocationID.Unlock()
			s.mxVisitsByUserID.Unlock()

			v.VisitedAt = visit.VisitedAt

			s.addVisitToVisitsByUserID(v)
			s.addVisitToVisitsByLocationID(v)
		}
	}
	return nil
}
