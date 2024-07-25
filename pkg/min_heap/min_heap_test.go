package min_heap

import (
	"container/heap"
	"testing"
)

func TestMinHeap_Len(t *testing.T) {
	h := &MinHeap{
		{Content: "1\n", Index: 0},
		{Content: "2\n", Index: 1},
		{Content: "3\n", Index: 2},
	}
	if h.Len() != 3 {
		t.Errorf("expected length 3, got %d", h.Len())
	}
}

func TestMinHeap_Less(t *testing.T) {
	h := &MinHeap{
		{Content: "1\n", Index: 0},
		{Content: "2\n", Index: 1},
	}
	if !h.Less(0, 1) {
		t.Errorf("expected h.Less(0, 1) to be true")
	}
	if h.Less(1, 0) {
		t.Errorf("expected h.Less(1, 0) to be false")
	}
}

func TestMinHeap_Swap(t *testing.T) {
	h := &MinHeap{
		{Content: "1\n", Index: 0},
		{Content: "2\n", Index: 1},
	}
	h.Swap(0, 1)
	if (*h)[0].Content != "2\n" || (*h)[1].Content != "1\n" {
		t.Errorf("expected swapped elements, got %v", h)
	}
}

func TestMinHeap_PushPop(t *testing.T) {
	h := &MinHeap{}
	heap.Init(h)

	heap.Push(h, LineFile{Content: "3\n", Index: 2})
	heap.Push(h, LineFile{Content: "1\n", Index: 0})
	heap.Push(h, LineFile{Content: "2\n", Index: 1})

	if h.Len() != 3 {
		t.Errorf("expected length 3, got %d", h.Len())
	}

	first := heap.Pop(h).(LineFile)
	if first.Content != "1\n" {
		t.Errorf("expected first popped element to be '1\\n', got %v", first.Content)
	}

	second := heap.Pop(h).(LineFile)
	if second.Content != "2\n" {
		t.Errorf("expected second popped element to be '2\\n', got %v", second.Content)
	}

	third := heap.Pop(h).(LineFile)
	if third.Content != "3\n" {
		t.Errorf("expected third popped element to be '3\\n', got %v", third.Content)
	}
}
