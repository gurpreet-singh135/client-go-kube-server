package service

import (
	"myapp/model"
	"container/heap"
	batchv1 "k8s.io/api/batch/v1"
	"sync"
)

type PriorityQueue []*model.CustomJob

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*model.CustomJob)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Update(item *model.CustomJob, value batchv1.Job, priority int) {
	item.Job = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}

type ThreadSafePriorityQueue struct {
	pq PriorityQueue
	mu sync.Mutex
}

// NewThreadSafePriorityQueue creates a new thread-safe priority queue.
func NewThreadSafePriorityQueue() *ThreadSafePriorityQueue {
	return &ThreadSafePriorityQueue{
		pq: make(PriorityQueue, 0),
	}
}

func (tspq *ThreadSafePriorityQueue) Push(x any) {

	tspq.mu.Lock()
	defer tspq.mu.Unlock()
	item := x.(*model.CustomJob)
	heap.Push(&tspq.pq, item)
}

// Pop removes and returns the element with the highest priority in a thread-safe manner.
func (tspq *ThreadSafePriorityQueue) Pop() (any, bool) {
	tspq.mu.Lock()
	defer tspq.mu.Unlock()
	if tspq.pq.Len() == 0 {
		return nil, false
	}
	item := heap.Pop(&tspq.pq)
	return item, true
}