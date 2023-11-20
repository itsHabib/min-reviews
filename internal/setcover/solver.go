package setcover

import (
	"fmt"
)

type Solver[T comparable] struct {
	bestSize   int
	chk        []T
	solutions  [][]T
	solTracker map[T]struct{}
	universe   Universe[T]
}

func NewSolver[T comparable](chk []T, universe Universe[T]) (*Solver[T], error) {
	if universe == nil {
		return nil, fmt.Errorf("unable to init new runner due to nil universe")
	}

	return &Solver[T]{chk: chk, universe: universe}, nil
}

func (s *Solver[T]) MinCover() [][]T {
	s.minCover([]T{}, 0)

	return s.solutions
}

func (s *Solver[T]) minCover(current []T, idx int) {
	if len(current) > s.bestSize && s.bestSize != 0 || idx == len(s.chk) {
		return
	}

	if s.universe.Covers(current) {
		if len(current) < s.bestSize {
			s.bestSize = len(current)
			s.solutions = [][]T{}
		}
		s.solutions = append(s.solutions, copySlice(current))
	}

	with := copySlice(current)
	with = append(with, s.chk[idx])
	s.minCover(current, idx+1)
	s.minCover(with, idx+1)
}

type Universe[T comparable] interface {
	Covers([]T) bool
}

func copySlice[T any](slice []T) []T {
	cp := make([]T, len(slice))
	copy(cp, slice)
	return cp
}
