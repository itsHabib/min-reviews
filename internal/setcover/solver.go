package setcover

import (
	"fmt"
	"sort"
	"strings"
)

type Solver struct {
	bestSize   int
	chk        []string
	solutions  [][]string
	solTracker map[string]struct{}
	universe   Universe
}

func NewSolver(chk []string, universe Universe) (*Solver, error) {
	if universe == nil {
		return nil, fmt.Errorf("unable to init new runner due to nil universe")
	}

	return &Solver{
		chk:        chk,
		universe:   universe,
		solTracker: make(map[string]struct{}),
	}, nil
}

func (s *Solver) MinCover() [][]string {
	s.minCover([]string{}, 0)

	return s.solutions
}

func (s *Solver) minCover(current []string, idx int) {
	if (len(current) > s.bestSize && s.bestSize != 0) || idx == len(s.chk) {
		return
	}

	if s.universe.Covers(current) {
		key := currentKey(current)
		if _, ok := s.solTracker[key]; ok {
			return
		}

		if s.bestSize == 0 || len(current) < s.bestSize {
			s.bestSize = len(current)
			s.solutions = [][]string{}
		}
		s.solutions = append(s.solutions, copySlice(current))
		s.solTracker[key] = struct{}{}
	}

	with := copySlice(current)
	with = append(with, s.chk[idx])
	s.minCover(current, idx+1)
	s.minCover(with, idx+1)
}

type Universe interface {
	Covers([]string) bool
}

func copySlice[string any](slice []string) []string {
	cp := make([]string, len(slice))
	copy(cp, slice)
	return cp
}

func currentKey(current []string) string {
	sort.Slice(current, func(i, j int) bool {
		return current[i] < current[j]
	})
	var b strings.Builder
	for i := range current {
		b.WriteString(current[i])
	}

	return b.String()
}
