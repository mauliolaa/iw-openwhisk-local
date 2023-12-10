package predictor

import (
	"container/heap"
	"fmt"
)

const (
	JAVA_PRIORITY       = 4
	RUBY_PRIORITY       = 3
	GO_PRIORITY         = 2
	JAVASCRIPT_PRIORITY = 2
	PYTHON_PRIORITY     = 1
	UNKNOWN_PRIORITY    = 10
)

func getPriorityValue(language string) int {
	switch language {
	case "js":
		return JAVASCRIPT_PRIORITY
	case "py":
		return PYTHON_PRIORITY
	case "go":
		return GO_PRIORITY
	case "rb":
		return RUBY_PRIORITY
	case "jar":
		return JAVA_PRIORITY
	default:
		return UNKNOWN_PRIORITY
	}
}

// PriorityFnRequest is something we manage in a priority queue.
type PriorityFnRequest struct {
	FnRequest
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*PriorityFnRequest

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*PriorityFnRequest)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Update(info map[string]any) {
	//TODO implement me
	item := &PriorityFnRequest{
		FnRequest: info["fnRequest"].(FnRequest),
		priority:  getPriorityValue(info["language"].(string)),
	}
	heap.Push(pq, item)
	print(pq.Len())
}

func (pq *PriorityQueue) Predict() *Prediction {
	if pq.Len() == 0 {
		fmt.Println("NILPREDICTION")
		return NilPrediction
	}
	item := heap.Pop(pq).(*PriorityFnRequest)
	fmt.Printf("%+v", item)
	return &Prediction{
		FnName:       item.FnName,
		FnParameters: item.FnParameters,
	}
}

func NewPriorityQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &pq
}
