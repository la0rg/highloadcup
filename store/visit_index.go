package store

import (
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
	return v.VisitedAt.Before(*(visit.VisitedAt))
}

func appendIterator(listPtr *[]model.Visit, country *string) func(item btree.Item) bool {
	return func(item btree.Item) bool {
		location := item.(VisitItem).Location
		if country != nil && (location == nil || *(location.Country) != *country) {
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

func (vi *VisitIndex) Get(fromDate *time.Time, toDate *time.Time, country *string) model.VisitArray {
	vi.mx.RLock()
	defer vi.mx.RUnlock()
	visits := make([]model.Visit, 0)
	switch {
	case fromDate == nil && toDate == nil:
		vi.byDate.Ascend(appendIterator(&visits, country))
	case fromDate != nil && toDate != nil:
		from := fromDate.Add(time.Microsecond)
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: &from}}
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vi.byDate.AscendRange(greaterOrEqual, lessThan, appendIterator(&visits, country))
	case fromDate != nil:
		from := fromDate.Add(time.Microsecond)
		greaterOrEqual := VisitItem{&model.Visit{VisitedAt: &from}}
		vi.byDate.AscendGreaterOrEqual(greaterOrEqual, appendIterator(&visits, country))
	case toDate != nil:
		lessThan := VisitItem{&model.Visit{VisitedAt: toDate}}
		vi.byDate.AscendLessThan(lessThan, appendIterator(&visits, country))
	}
	return model.VisitArray{
		Visits: visits,
	}
}
