// Copyright © 2026 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package rendezvous

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestStack_Basic(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New[int](ctx)

	if !s.IsEmpty() {
		t.Errorf("IsEmpty() got = %v, want %v", false, true)
	}

	if s.Size() != 0 {
		t.Errorf("Size() got = %d, want %d", s.Size(), 0)
	}

	s.Push(10)

	if s.IsEmpty() {
		t.Errorf("IsEmpty() got = %v, want %v", true, false)
	}

	if s.Size() != 1 {
		t.Errorf("Size() got = %d, want %d", s.Size(), 1)
	}

	val, ok := s.Peek()
	if !ok || val != 10 {
		t.Errorf("Peek() got = (%v, %v), want = (%v, %v)", val, ok, 10, true)
	}

	if s.Size() != 1 {
		t.Errorf("Size() after Peek() got = %d, want %d", s.Size(), 1)
	}

	val, ok = s.Pop()
	if !ok || val != 10 {
		t.Errorf("Pop() got = (%v, %v), want = (%v, %v)", val, ok, 10, true)
	}

	if s.Size() != 0 {
		t.Errorf("Size() after Pop() got = %d, want %d", s.Size(), 0)
	}
}

func TestStack_Empty(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New[string](ctx)

	val, ok := s.Peek()
	if ok || val != "" {
		t.Errorf("Peek() on empty stack got = (%v, %v), want = (%v, %v)", val, ok, "", false)
	}

	val, ok = s.Pop()
	if ok || val != "" {
		t.Errorf("Pop() on empty stack got = (%v, %v), want = (%v, %v)", val, ok, "", false)
	}
}

func TestStack_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := New[int](ctx)

	s.Push(1)
	cancel() // Cancel the context, server goroutine should stop.

	// Give the scheduler a moment to stop the goroutine.
	time.Sleep(50 * time.Millisecond)

	done := make(chan struct{})
	go func() {
		// This should block forever as the server is gone.
		s.Size()
		close(done)
	}()

	select {
	case <-done:
		t.Error("s.Size() returned, but should have blocked forever after context cancellation")
	case <-time.After(100 * time.Millisecond):
		// This is the expected outcome. The call to Size() blocked.
	}
}

func TestStack_ConcurrentPushPop(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New[int](ctx)
	numGoroutines := 50
	numItems := 100
	var wg sync.WaitGroup

	// Concurrent Pushes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numItems; j++ {
				s.Push(j)
			}
		}()
	}
	wg.Wait()

	expectedSize := numGoroutines * numItems
	if s.Size() != expectedSize {
		t.Errorf("Size() after concurrent pushes got = %d, want %d", s.Size(), expectedSize)
	}
}
