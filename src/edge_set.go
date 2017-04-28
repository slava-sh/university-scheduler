package main

import (
	"container/heap"
	"math/rand"
)

type EdgeSet struct {
	pq    PriorityQueue
	items map[Edge]*Item
}

type PriorityQueue []*Item

type Edge struct {
	Group int
	Prof  int
	Day   int
	Class int
}

type Item struct {
	value    Edge
	priority int
	index    int
}

func NewEdgeSet() *EdgeSet {
	return &EdgeSet{
		make(PriorityQueue, 0),
		make(map[Edge]*Item),
	}
}

func (es *EdgeSet) Push(g, p, d, c int) {
	edge := Edge{g, p, d, c}
	if es.items[edge] != nil {
		return
	}
	item := &Item{
		value:    edge,
		priority: rand.Int(),
	}
	heap.Push(&es.pq, item)
	es.items[edge] = item
}

func (es *EdgeSet) Pop() (g, p, d, c int) {
	edge := heap.Pop(&es.pq).(*Item).value
	delete(es.items, edge)
	return edge.Group, edge.Prof, edge.Day, edge.Class
}

func (es *EdgeSet) Remove(g, p, d, c int) {
	edge := Edge{g, p, d, c}
	item := es.items[edge]
	if item == nil {
		return
	}
	heap.Remove(&es.pq, item.index)
	delete(es.items, edge)
}

func (es *EdgeSet) Len() int {
	return es.pq.Len()
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority >= pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // For safety.
	*pq = old[:n-1]
	return item
}
