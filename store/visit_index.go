package store

import (
	"github.com/google/btree"
	"github.com/la0rg/highloadcup/model"
	"github.com/mailru/easyjson/opt"
)

type VisitIndex struct {
	//mx     sync.RWMutex
	byDate *btree.BTree
}

type VisitItem struct {
	*model.Visit
}

func (v VisitItem) Less(then btree.Item) bool {
	// VisitedAt as a main index of the tree
	// TODO; may be bug in deletion if we have visits with the same timestamp
	return v.VisitedAt.V < then.(VisitItem).VisitedAt.V
}

func appendIteratorByCountryAndVisit(listPtr *[]model.Visit, country *string, toDistance *int32) func(item btree.Item) bool {
	return func(item btree.Item) bool {
		location := item.(VisitItem).Location
		// country - название страны, в которой находятся интересующие достопримечательности
		if country != nil && (location == nil || location.Country.V != *country) {
			return true
		}
		// toDistance - возвращать только те места, у которых расстояние от города меньше этого параметра
		if toDistance != nil && (location == nil || location.Distance.V >= *toDistance) {
			return true
		}

		*listPtr = append(*listPtr, *(item.(VisitItem).Visit))
		return true
	}
}

func appendIteratorByAgeAndGender(listPtr *[]model.Visit, fromAge *int64, toAge *int64, gender *string) func(item btree.Item) bool {
	return func(item btree.Item) bool {
		user := item.(VisitItem).User
		// fromAge - учитывать только путешественников, у которых возраст (считается от текущего timestamp) строго больше этого параметра
		// birthdate < timestamp
		if fromAge != nil && (user == nil || *fromAge <= user.BirthDate.V) {
			return true
		}
		// birthdate > timestamp
		if toAge != nil && (user == nil || user.BirthDate.V <= *toAge) {
			return true
		}
		if gender != nil && (user == nil || user.Gender.V != *gender) {
			return true
		}
		*listPtr = append(*listPtr, *(item.(VisitItem).Visit))
		return true
	}
}

func NewVisitIndex() *VisitIndex {
	return &VisitIndex{
		byDate: btree.New(5),
	}
}

func (vi *VisitIndex) Add(visit *model.Visit) {
	//vi.mx.Lock()
	//defer vi.mx.Unlock()
	vi.byDate.ReplaceOrInsert(VisitItem{visit})
}

func (vi *VisitIndex) get(fromDate *int64, toDate *int64, iter btree.ItemIterator) {
	switch {
	case fromDate == nil && toDate == nil:
		vi.byDate.Ascend(iter)
	case fromDate != nil && toDate != nil:
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: opt.Int64{V: *fromDate + 1, Defined: true}}}
		lessThan := VisitItem{&model.Visit{VisitedAt: opt.Int64{V: *toDate, Defined: true}}}
		vi.byDate.AscendRange(greaterOrEqual, lessThan, iter)
	case fromDate != nil:
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: opt.Int64{V: *fromDate + 1, Defined: true}}}
		vi.byDate.AscendGreaterOrEqual(greaterOrEqual, iter)
	case toDate != nil:
		lessThan := VisitItem{&model.Visit{VisitedAt: opt.Int64{V: *toDate, Defined: true}}}
		vi.byDate.AscendLessThan(lessThan, iter)
	}
}

func (vi *VisitIndex) GetByCountryAndDistance(fromDate *int64, toDate *int64, country *string, toDistance *int32) model.UserVisitArray {
	//vi.mx.RLock()
	//defer vi.mx.RUnlock()
	visits := make([]model.Visit, 0)
	vi.get(fromDate, toDate, appendIteratorByCountryAndVisit(&visits, country, toDistance))
	userVisits := make([]model.UserVisit, len(visits))
	for i := range visits {
		userVisits[i] = model.UserVisitFromVisit(visits[i])
	}
	return model.UserVisitArray{
		Visits: userVisits,
	}
}

func (vi *VisitIndex) GetByAgeAndGender(fromDate *int64, toDate *int64, fromAge *int64, toAge *int64, gender *string) []model.Visit {
	//vi.mx.RLock()
	//defer vi.mx.RUnlock()
	visits := make([]model.Visit, 0)
	vi.get(fromDate, toDate, appendIteratorByAgeAndGender(&visits, fromAge, toAge, gender))
	return visits
}

func (vi *VisitIndex) Remove(visit *model.Visit) bool {
	//vi.mx.Lock()
	//defer vi.mx.Unlock()
	deleted := vi.byDate.Delete(VisitItem{visit})
	return deleted != nil
}

func (vi *VisitIndex) ApplyToAll(f func(*model.Visit)) {
	//vi.mx.Lock()
	//defer vi.mx.Unlock()
	vi.byDate.Ascend(func(item btree.Item) bool {
		f(item.(VisitItem).Visit)
		return true
	})
}
