package util

import (
	"math/rand"
	"testing"
)

func TestQueue(t *testing.T) {
	testCases := []struct {
		name          string
		initialQueue  []string
		enqueueItems  []string
		dequeCount    int
		expectedItems []string
		expectedSize  int
		dequeError    bool
		clearQueue    bool
	}{
		{
			name:          "Basic Enqueue and Dequeue",
			enqueueItems:  []string{"http://example.com/1", "http://example.com/2", "http://example.com/3"},
			dequeCount:    3,
			expectedItems: []string{"http://example.com/1", "http://example.com/2", "http://example.com/3"},
			expectedSize:  0,
			dequeError:    true,
		},
		{
			name:          "Enqueue and Size",
			enqueueItems:  []string{"http://example.com/1", "http://example.com/2", "http://example.com/3"},
			expectedSize:  3,
			dequeCount:    0,
			expectedItems: []string{},
			dequeError:    false,
		},
		{
			name:          "Dequeue from Empty Queue",
			expectedSize:  0,
			dequeCount:    1,
			expectedItems: []string{},
			dequeError:    true,
		},
		{
			name:          "Clear Queue",
			enqueueItems:  []string{"http://example.com/1", "http://example.com/2", "http://example.com/3"},
			clearQueue:    true,
			expectedSize:  0,
			dequeCount:    0,
			expectedItems: []string{},
			dequeError:    false,
		},
		{
			name:          "Initial Queue",
			initialQueue:  []string{"http://example.com/1", "http://example.com/2"},
			dequeCount:    2,
			expectedItems: []string{"http://example.com/1", "http://example.com/2"},
			expectedSize:  0,
			dequeError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewQueue()

			// Initialize queue with initial items
			for _, item := range tc.initialQueue {
				q.Enque(item)
			}

			// Enqueue items
			for _, item := range tc.enqueueItems {
				q.Enque(item)
			}

			// Check size after enqueueing, if a size is specified
			if tc.expectedSize != 0 {
				if size := q.Size(); size != tc.expectedSize {
					t.Errorf("Expected size %d, got %d", tc.expectedSize, size)
				}
			}

			// Dequeue items and check their values
			for i := 0; i < tc.dequeCount; i++ {
				item, err := q.Deque()
				if i < len(tc.expectedItems) {
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					if item != tc.expectedItems[i] {
						t.Errorf("Expected item '%s', got '%s'", tc.expectedItems[i], item)
					}
				} else {
					if tc.dequeError {
						if err == nil {
							t.Error("Expected error when dequeuing from empty queue, got nil")
						}
					} else if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
				}
			}

			// Clear queue
			if tc.clearQueue {
				q.Clear()
			}

			// Verify the queue is empty after clearing
			if tc.clearQueue && len(q.items) != 0 {
				t.Error("Expected empty queue, got a queue with items")
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	testCases := []struct {
		name         string
		initialQueue []string
		seed         int64
	}{
		{
			name:         "Basic Shuffle",
			initialQueue: []string{"A", "B", "C"},
			seed:         121,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewQueue()
			for _, item := range tc.initialQueue {
				q.Enque(item)
			}

			originalOrder := make([]string, len(q.items))
			copy(originalOrder, q.items)

			rand.Seed(tc.seed)
			q.Shuffle()

			if len(q.items) != len(originalOrder) {
				t.Error("Shuffled queue length does not match original length")
			}

			isDifferent := false
			for i := range originalOrder {
				if q.items[i] != originalOrder[i] {
					isDifferent = true
					break
				}
			}

			if !isDifferent {
				t.Error("Shuffled queue order should be different from original order")
			}
		})
	}
}
