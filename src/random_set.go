package main

import (
	"container/heap"
	"math/rand"
)

type RandomSet PriorityQueue
type PriorityQueue []*Item

type Item struct {
	value    *Solution
	priority int
}

func MakeRandomSet() RandomSet {
	return make(RandomSet, 0)
}

func (rs *RandomSet) Push(s *Solution) {
	pq := PriorityQueue(*rs)
	heap.Push(&pq, &Item{s, rand.Int()})
	*rs = RandomSet(pq)
}

func (rs *RandomSet) Pop() *Solution {
	pq := PriorityQueue(*rs)
	solution := heap.Pop(&pq).(*Item).value
	*rs = RandomSet(pq)
	return solution
}

func (rs RandomSet) Len() int {
	pq := PriorityQueue(rs)
	return pq.Len()
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority >= pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
