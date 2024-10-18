package main

type IntSet struct {
	data map[int]struct{}

	/* Item count */
	Count int
}

func NewIntSet() (s IntSet) {
	s.data = make(map[int]struct{})

	return
}

func (s *IntSet) Exists(n int) (ok bool) {

	_, ok = s.data[n]

	return
}

func (s *IntSet) Insert(n int) bool {
	if s.Exists(n) {

		return false
	}

	s.data[n] = struct{}{}

	s.Count++

	return true
}
