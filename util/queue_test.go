package util

import "testing"

func TestQueue(t *testing.T) {
    q := NewQueue()

    q.Enque("http://example.com/1")
    q.Enque("http://example.com/2")
    q.Enque("http://example.com/3")

    if size := q.Size(); size != 3 {
        t.Errorf("Expected size 3, got %d", size)
    }

    item, err := q.Deque()
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if item != "http://example.com/3" {
        t.Errorf("Expected item 'http://example.com/3', got '%s'", item)
    }

    if size := q.Size(); size != 2 {
        t.Errorf("Expected size 2, got %d", size)
    }

    item, err = q.Deque()
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    
    if item != "http://example.com/2" {
        t.Errorf("Expected item 'http://example.com/2', got '%s'", item)
    }

    item, err = q.Deque()
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    
    if item != "http://example.com/1" {
        t.Errorf("Expected item 'http://example.com/1', got '%s'", item)
    }

    _, err = q.Deque()
    if err == nil {
        t.Error("Expected error when dequeuing from empty queue, got nil")
    }
	

    q.Enque("http://example.com/1")
    q.Enque("http://example.com/2")
    q.Enque("http://example.com/3")
	
    q.Clear()

    if len(q.items) !=0 {
	t.Error("Expected empty queue, got a queue with items")
    }
}

func TestShuffle(t *testing.T) {
	q := NewQueue()
	q.Enque("A")
	q.Enque("B")
	q.Enque("C")

	originalOrder := make([]string, len(q.items))
	copy(originalOrder, q.items)

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
}
