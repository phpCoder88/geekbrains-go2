package main

import "sync"

type SetRW struct {
	sync.RWMutex
	mm map[int]struct{}
}

func NewSetRW() *SetRW {
	return &SetRW{
		mm: map[int]struct{}{},
	}
}

func (s *SetRW) Add(i int) {
	s.Lock()
	s.mm[i] = struct{}{}
	s.Unlock()
}

func (s *SetRW) Has(i int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.mm[i]
	return ok
}
