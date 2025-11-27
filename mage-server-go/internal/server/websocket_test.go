package server

import (
	"sync"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/session"
)

// TestWebSocketCloseRaceCondition tests that replacing a WebSocket connection
// doesn't cause a "close of closed channel" panic when the old connection's
// goroutines try to close their done channel after it's been closed externally.
func TestWebSocketCloseRaceCondition(t *testing.T) {
	// This test simulates the race condition that occurs when:
	// 1. Connection A is active with its own done channel and closeDone function
	// 2. Connection B comes in and closes Connection A's channel directly
	// 3. Connection A's goroutines then try to call their closeDone function
	//
	// The bug: Connection B uses close(oldCloseChan) directly, bypassing
	// Connection A's sync.Once, so when Connection A's goroutines call
	// closeDone(), the sync.Once.Do() executes close() on an already-closed channel.

	t.Run("direct close bypasses sync.Once causing panic (demonstrates bug pattern)", func(t *testing.T) {
		// This test demonstrates the BUG PATTERN that we fixed.
		// It shows what happens if you use close() directly instead of the closeDone function.

		// Simulate Connection A's setup
		doneA := make(chan struct{})
		var closeOnceA sync.Once
		closeDoneA := func() {
			closeOnceA.Do(func() {
				close(doneA)
			})
		}

		// Simulate the BUGGY pattern: closing channel directly instead of using closeDone
		close(doneA)

		// Now simulate Connection A's readHandler calling closeDone after read error
		// This demonstrates WHY the bug happens
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic (demonstrates bug pattern): %v", r)
				// This is the bug we fixed by storing closeDone in the session
			}
		}()

		closeDoneA()   // This panics because channel was closed directly!
		_ = closeDoneA // silence unused warning
	})

	t.Run("using closeDone function prevents panic", func(t *testing.T) {
		// This test shows the correct behavior after the fix:
		// Store and call the closeDone function instead of closing the channel directly

		doneA := make(chan struct{})
		var closeOnceA sync.Once
		closeDoneA := func() {
			closeOnceA.Do(func() {
				close(doneA)
			})
		}

		// Connection B should call the closeDone function, not close() directly
		closeDoneA()

		// Now when Connection A's goroutines call closeDone, it's safe
		// because sync.Once ensures it only runs once
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("UNEXPECTED PANIC: %v - the fix should prevent this", r)
			}
		}()

		closeDoneA() // This should NOT panic
		closeDoneA() // Can be called multiple times safely
		closeDoneA()
	})
}

// TestWebSocketConnectionReplacement tests the full flow of replacing connections
func TestWebSocketConnectionReplacement(t *testing.T) {
	t.Run("concurrent goroutines closing safely", func(t *testing.T) {
		done := make(chan struct{})
		var closeOnce sync.Once
		closeDone := func() {
			closeOnce.Do(func() {
				close(done)
			})
		}

		// Simulate multiple goroutines trying to close concurrently
		// (readHandler, pingHandler, main loop, and external close)
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panic in goroutine: %v", r)
					}
				}()
				closeDone()
			}()
		}

		wg.Wait()

		// Verify channel is closed
		select {
		case <-done:
			// Good, channel is closed
		default:
			t.Error("done channel should be closed")
		}
	})

	t.Run("external close followed by internal close", func(t *testing.T) {
		// Simulates the exact scenario from the bug:
		// 1. Connection A is running
		// 2. Connection B closes A's channel
		// 3. A's readHandler gets error and tries to close

		doneA := make(chan struct{})
		var closeOnceA sync.Once
		closeDoneA := func() {
			closeOnceA.Do(func() {
				close(doneA)
			})
		}

		// Simulate readHandler waiting on channel
		readHandlerDone := make(chan struct{})
		go func() {
			defer close(readHandlerDone)
			// Simulate blocking read that will unblock when done is closed
			<-doneA
			// After unblocking, handler tries to close (simulating error path)
			closeDoneA()
		}()

		// Simulate external close (from new connection)
		time.Sleep(10 * time.Millisecond)
		closeDoneA()

		// Wait for readHandler to finish
		select {
		case <-readHandlerDone:
			// Success - no panic
		case <-time.After(time.Second):
			t.Error("readHandler timed out")
		}
	})
}

// TestSessionWebSocketCloseFunc tests the actual session integration for connection replacement
func TestSessionWebSocketCloseFunc(t *testing.T) {
	t.Run("connection replacement using session close function", func(t *testing.T) {
		// Create a session
		sess := session.NewSession("test-session", "localhost", time.Hour)

		// Simulate Connection A
		doneA := make(chan struct{})
		var closeOnceA sync.Once
		closeDoneA := func() {
			closeOnceA.Do(func() {
				close(doneA)
			})
		}

		// Register Connection A with session
		oldFunc := sess.SetWebSocketCloseFunc(doneA, closeDoneA)
		if oldFunc != nil {
			t.Error("expected nil old function for first connection")
		}

		// Start a goroutine simulating Connection A's readHandler
		readHandlerDone := make(chan struct{})
		go func() {
			defer close(readHandlerDone)
			<-doneA // Wait for close signal
			// Simulate readHandler trying to close after read error
			closeDoneA()
		}()

		// Simulate Connection B coming in
		doneB := make(chan struct{})
		var closeOnceB sync.Once
		closeDoneB := func() {
			closeOnceB.Do(func() {
				close(doneB)
			})
		}

		// Register Connection B - this should return A's close function
		oldFunc = sess.SetWebSocketCloseFunc(doneB, closeDoneB)
		if oldFunc == nil {
			t.Fatal("expected old close function from Connection A")
		}

		// Use the returned close function to safely close Connection A
		// This is the KEY FIX: using the function instead of close() directly
		oldFunc()

		// Wait for Connection A's readHandler to finish
		select {
		case <-readHandlerDone:
			// Success - no panic!
		case <-time.After(time.Second):
			t.Error("readHandler timed out")
		}

		// Verify Connection A is closed
		select {
		case <-doneA:
			// Good
		default:
			t.Error("Connection A's done channel should be closed")
		}
	})

	t.Run("rapid connection replacement", func(t *testing.T) {
		sess := session.NewSession("test-session", "localhost", time.Hour)

		// Rapidly replace connections - simulates user refreshing page repeatedly
		for i := 0; i < 100; i++ {
			done := make(chan struct{})
			var closeOnce sync.Once
			closeDone := func() {
				closeOnce.Do(func() {
					close(done)
				})
			}

			// Start a goroutine that will try to close when signaled
			var wg sync.WaitGroup
			wg.Add(1)
			go func(closeFn func()) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panic during iteration: %v", r)
					}
				}()
				// Simulate some work then try to close
				time.Sleep(time.Microsecond)
				closeFn()
			}(closeDone)

			// Register with session, close previous
			oldFunc := sess.SetWebSocketCloseFunc(done, closeDone)
			if oldFunc != nil {
				oldFunc() // Safely close previous connection
			}

			wg.Wait()
		}
	})
}
