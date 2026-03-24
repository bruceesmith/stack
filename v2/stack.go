/*
Package stack defines goroutine-safe methods for manipulating a generic [stack] data structure via the standard operations
IsEmpty, Peek, Pop, Push and Size. Sub-folders hold differing implementations of the Stack - mutex/ uses locking via the
standard sync package; rendezvous/ uses a synchronous actor implementation via channels.

Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
Use of this source code is governed by the MIT
License that can be found in the LICENSE file.

[stack]: https://en.wikipedia.org/wiki/Stack_(abstract_data_type)
*/
package stack

//go:generate ./make_doc.sh

// Stack represents a thread-safe LIFO (Last-In-First-Out) collection.
type Stack[T any] interface {
	// IsEmpty checks if the stack has no elements.
	IsEmpty() bool

	// Peek returns the top element without removing it.
	// It returns false if the stack is empty.
	Peek() (value T, ok bool)

	// Pop removes and returns the top element.
	// It returns false if the stack is empty.
	Pop() (value T, ok bool)

	// Push adds an element to the top of the stack.
	Push(value T)

	// Size returns the number of elements in the stack.
	Size() int
}
