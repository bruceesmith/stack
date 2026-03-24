// Copyright © 2026 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package rendezvous defines Stack implementation uses a synchronous actor model via channels
*/
package rendezvous

import (
	"context"
)

// element is an entry on the stack
type element[T any] struct {
	value T
	next  *element[T]
}

// response is a response from the actor (server)
type response[T any] struct {
	i     int
	ok    bool
	value T
}

// action is the type of action requested by ask()
type action = uint8

const (
	empty action = iota
	peek
	pop
	push
	size
)

// message is an incoming request to ask()
type message[T any] struct {
	action action
	value  T
}

// Stack is a Go stack implementation using a linked list
// It is go-routine safe
type Stack[T any] struct {
	inbox  chan message[T]
	outbox chan response[T]
	size   int
	top    *element[T]
	zero   T
}

// New creates an empty Stack and starts the actor goroutine
func New[T any](ctx context.Context) *Stack[T] {
	s := &Stack[T]{
		inbox:  make(chan message[T]),
		outbox: make(chan response[T]),
		zero:   *new(T),
	}
	go s.server(ctx, s.inbox)
	return s
}

// ask posts a request, waits for the response, returns that response
func (s *Stack[T]) ask(command action, value T) response[T] {
	s.inbox <- message[T]{
		action: command,
		value:  value,
	}
	r := <-s.outbox
	return r
}

// IsEmpty returns true if the stack has no elements
func (s *Stack[T]) IsEmpty() bool {
	r := s.ask(empty, s.zero)
	return r.ok
}

func (s *Stack[T]) handleIsEmpty() response[T] {
	return response[T]{
		ok: s.size == 0,
	}
}

// Peek returns a copy of the top element off the stack
func (s *Stack[T]) Peek() (value T, ok bool) {
	r := s.ask(peek, s.zero)
	return r.value, r.ok

}

func (s *Stack[T]) handlePeek() response[T] {
	var v T
	if s.top == nil {
		return response[T]{
			ok:    false,
			value: v,
		}
	}
	v = s.top.value
	return response[T]{
		ok:    true,
		value: v,
	}

}

// Pop removes the top element and returns it
func (s *Stack[T]) Pop() (value T, ok bool) {
	r := s.ask(pop, s.zero)
	return r.value, r.ok

}

func (s *Stack[T]) handlePop() response[T] {
	var v T
	if s.top == nil {
		return response[T]{
			ok:    false,
			value: v,
		}
	}
	v = s.top.value
	s.top = s.top.next
	s.size--
	return response[T]{
		ok:    true,
		value: v,
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(v T) {
	s.ask(push, v)
}

func (s *Stack[T]) handlePush(v T) {
	element := &element[T]{
		value: v,
		next:  s.top,
	}
	s.top = element
	s.size++
}

// server is the heart of the actor implementation. Callers use ask()
// to rendezvous with the server, thus single threading requests without
// the use of locks.
func (s *Stack[T]) server(ctx context.Context, inbox <-chan message[T]) {
loop:
	for {
		select {
		case msg := <-inbox:
			switch msg.action {
			case empty:
				s.outbox <- s.handleIsEmpty()
			case peek:
				s.outbox <- s.handlePeek()
			case pop:
				s.outbox <- s.handlePop()
			case push:
				s.handlePush(msg.value)
				s.outbox <- response[T]{}
			case size:
				s.outbox <- s.handleSize()
			}
		case <-ctx.Done():
			break loop
		}
	}
}

// Size returns the number of elements on the stack
func (s *Stack[T]) Size() int {
	r := s.ask(size, s.zero)
	return r.i
}

func (s *Stack[T]) handleSize() response[T] {
	return response[T]{
		i: s.size,
	}
}
