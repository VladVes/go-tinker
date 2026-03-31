package basics

import (
	"fmt"
	"slices"
)

type SetStorage interface {
	comparable
}

type Set[T SetStorage] struct {
	values []T
}

func NewSet[T SetStorage]() Set[T] {
	return Set[T]{}
}

func (s *Set[T]) Add(v T) {
	s.values = append(s.values, v)
}

func (s *Set[T]) Has(v T) bool {
	return slices.Contains(s.values, v)
}

func (s *Set[T]) Remove(v T) {
	result := s.values[:0] // переиспользование памяти
	for _, item := range s.values {
		if item != v {
			fmt.Println(item)
			result = append(result, item)
		}
	}
	s.values = result
}
