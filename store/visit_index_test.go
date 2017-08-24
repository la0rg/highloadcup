package store

import (
	"testing"

	"github.com/la0rg/highloadcup/model"
)

func TestVisitIndex_ApplyToAll(t *testing.T) {
	vi := NewVisitIndex()
	location := &model.Location{}
	var t1, t2, t3 int64
	t1 += 1
	t2 += 2
	t3 += 3
	visit1 := model.Visit{VisitedAt: &t1}
	visit2 := model.Visit{VisitedAt: &t2}
	visit3 := model.Visit{VisitedAt: &t3}
	vi.Add(&visit1)
	vi.Add(&visit2)
	vi.Add(&visit3)

	vi.ApplyToAll(func(visit *model.Visit) {
		visit.Location = location
	})

	if visit1.Location != location || visit2.Location != location || visit3.Location != location {
		t.Errorf("visit1.Location: %p, location: %p", visit1.Location, location)
	}
}
