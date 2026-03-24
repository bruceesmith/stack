package stack_test

import (
	"context"
	"fmt"

	"github.com/bruceesmith/stack/v2"
	"github.com/bruceesmith/stack/v2/mutex"
	"github.com/bruceesmith/stack/v2/rendezvous"
)

func Example() {
	var s1 stack.Stack[int] = mutex.New[int]()
	var s2 stack.Stack[int] = rendezvous.New[int](context.Background())

	s1.Push(1)
	s2.Push(2)
	fmt.Println("S1 size", s1.Size(), "s2 size", s2.Size())
}
