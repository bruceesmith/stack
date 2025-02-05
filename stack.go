// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package stack defines goroutine-safe methods for manipulating a generic [stack] data structure via the standard operations IsEmpty,	Peek, Pop, Push and Size.

[stack]: https://en.wikipedia.org/wiki/Stack_(abstract_data_type)
*/
package stack

//go:generate ./make_doc.sh

import "sync"

// element is an entry on the stack
type element[T any] struct {
	value T
	next  *element[T]
}

// Stack is a Go stack implementation using a linked  list
// It is go-routine safe
type Stack[T any] struct {
	top  *element[T]
	mut  *sync.Mutex
	size int
}

// New creates an empty Stack
func New[T any]() *Stack[T] {
	return &Stack[T]{
		mut: &sync.Mutex{},
	}
}

// IsEmpty returns true if the stack has no elements
func (s *Stack[T]) IsEmpty() bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.size == 0
}

// Peek returns a copy of the top element off the stack
func (s *Stack[T]) Peek() (value T, ok bool) {
	s.mut.Lock()
	defer s.mut.Unlock()

	var v T
	if s.top == nil {
		return v, false
	}

	v = s.top.value
	return v, true
}

// Pop removes the top element and returns it
func (s *Stack[T]) Pop() (value T, ok bool) {
	s.mut.Lock()
	defer s.mut.Unlock()

	var v T
	if s.top == nil {
		return v, false
	}

	v = s.top.value
	s.top = s.top.next
	s.size--

	return v, true
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(v T) {
	s.mut.Lock()
	defer s.mut.Unlock()

	element := &element[T]{
		value: v,
		next:  s.top,
	}

	s.top = element
	s.size++
}

// Size returns the number of elements on the stack
func (s *Stack[T]) Size() int {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.size
}
