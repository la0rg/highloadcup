package store

import (
	"github.com/la0rg/highloadcup/model"
)

type VisitList struct {
	m map[int32]*model.Visit
}

func NewVisitList() VisitList {
	return VisitList{make(map[int32]*model.Visit)}
}

func (vl *VisitList) Add(visit *model.Visit) {
	vl.m[*(visit.ID)] = visit
}

func (vl *VisitList) Remove(visit *model.Visit) {
	delete(vl.m, *(visit.ID))
}

func (vl *VisitList) Iter() chan *model.Visit {
	c := make(chan *model.Visit)
	go func() {
		for k := range vl.m {
			c <- vl.m[k]
		}
		close(c)
	}()
	return c
}
