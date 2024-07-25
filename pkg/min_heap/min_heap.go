package min_heap

import (
	"ext-sort/pkg/converters"
)

type MinHeap []LineFile

type LineFile struct {
	Content string
	Index   int
}

func (h MinHeap) Len() int { return len(h) }

func (h MinHeap) Less(i, j int) bool {
	return converters.StringAsInt(h[i].Content) < converters.StringAsInt(h[j].Content)
}

func (h MinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(LineFile))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
