package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/btree"
	"github.com/la0rg/highloadcup/model"
)

type VisitIndex struct {
	mx     sync.RWMutex
	byDate *btree.BTree
}

type VisitItem struct {
	*model.Visit
}

func (v VisitItem) Less(then btree.Item) bool {
	visit := then.(VisitItem)
	//log.Printf("Compare %v %v %p and %v %v %p", (*(v.VisitedAt)).Unix(), *(v.ID), v.Visit, (*(visit.VisitedAt)).Unix(), *(visit.ID), visit.Visit)
	// VisitedAt as a main index of the tree
	return v.VisitedAt.Before(*(visit.VisitedAt)) //|| (visit.ID != nil && *(v.ID) < *(visit.ID))
}

func appendIterator(listPtr *[]model.Visit, country *string, toDistance *int32) func(item btree.Item) bool {
	return func(item btree.Item) bool {
		location := item.(VisitItem).Location
		// country - название страны, в которой находятся интересующие достопримечательности
		if country != nil && (location == nil || *(location.Country) != *country) {
			return true
		}
		// toDistance - возвращать только те места, у которых расстояние от города меньше этого параметра
		if toDistance != nil && (location == nil || *(location.Distance) >= *toDistance) {
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
	vi.mx.Lock()
	defer vi.mx.Unlock()
	vi.byDate.ReplaceOrInsert(VisitItem{visit})
}

func (vi *VisitIndex) Get(fromDate *time.Time, toDate *time.Time, country *string, toDistance *int32) model.UserVisitArray {
	vi.mx.RLock()
	defer vi.mx.RUnlock()
	visits := make([]model.Visit, 0)
	switch {
	case fromDate == nil && toDate == nil:
		vi.byDate.Ascend(appendIterator(&visits, country, toDistance))
	case fromDate != nil && toDate != nil:
		from := fromDate.Add(time.Microsecond)
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: &from}}
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vi.byDate.AscendRange(greaterOrEqual, lessThan, appendIterator(&visits, country, toDistance))
	case fromDate != nil:
		from := fromDate.Add(time.Microsecond)
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: &from}}
		vi.byDate.AscendGreaterOrEqual(greaterOrEqual, appendIterator(&visits, country, toDistance))
	case toDate != nil:
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vi.byDate.AscendLessThan(lessThan, appendIterator(&visits, country, toDistance))
	}
	userVisits := make([]model.UserVisit, len(visits))
	for i := range visits {
		userVisits[i] = model.UserVisit{Visit: visits[i]}
	}
	return model.UserVisitArray{
		Visits: userVisits,
	}
}

func (vi *VisitIndex) Remove(visit *model.Visit) bool {
	vi.mx.Lock()
	defer vi.mx.Unlock()
	deleted := vi.byDate.Delete(VisitItem{visit})
	fmt.Println(deleted)
	return deleted != nil
}
