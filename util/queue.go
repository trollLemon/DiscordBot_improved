package util

import "sync"
import "fmt"
import "math/rand"

type Queue struct {
	items []string
	lock  sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]string, 0),
	}
}

func (q *Queue) _Swap(i int, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}
func (q *Queue) Enque(url string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = append(q.items, url)
}
func (q *Queue) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.items)
}
func (q *Queue) Deque() (string, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	items := len(q.items)

	if items == 0 {
		return "", fmt.Errorf("Empty Queue")
	}

	first := q.items[0]
	q.items = q.items[1:items]

	return first, nil

}

func (q *Queue) Clear() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = nil
	q.items = make([]string,0)
}


func (q *Queue) Shuffle() {

	for i := len(q.items) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		q._Swap(i, j)
	}
}
