package main

import "math/rand"

// Map-like persistent treap.
type Treap struct {
	root *node
}

type key3 struct {
	a int
	b int
	c int
}

type node struct {
	key      key3
	minC     int
	maxC     int
	value    int
	priority int
	size     int
	left     *node
	right    *node
}

func newNode(key key3, value int) *node {
	return &node{
		key:      key,
		minC:     key.c,
		maxC:     key.c,
		value:    value,
		size:     1,
		priority: rand.Int(),
	}
}

func (n *node) copy() *node {
	if n == nil {
		return nil
	}
	copy := new(node)
	*copy = *n
	return copy
}

func (n *node) update() {
	if n == nil {
		return
	}
	n.size = 1
	n.minC = n.key.c
	n.maxC = n.key.c
	if n.left != nil {
		n.size += n.left.size
		n.minC = min(n.minC, n.left.minC)
		n.maxC = max(n.maxC, n.left.maxC)
	}
	if n.right != nil {
		n.size += n.right.size
		n.minC = min(n.minC, n.right.minC)
		n.maxC = max(n.maxC, n.right.maxC)
	}
}

func merge(left, right *node) (result *node) {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.priority >= right.priority {
		result = left.copy()
		result.right = merge(result.right, right)
	} else {
		result = right.copy()
		result.left = merge(left, result.left)
	}
	result.update()
	return
}

func split2(node *node, key key3) (left *node, right *node) {
	if node == nil {
		return
	}
	if cmp(node.key, key) <= 0 {
		left = node.copy()
		left.right, right = split2(node.right, key)
		left.update()
	} else {
		right = node.copy()
		left, right.left = split2(node.left, key)
		right.update()
	}
	return
}

func split3(node *node, key key3) (left, middle, right *node) {
	prevKey := key3{key.a, key.b, key.c - 1}
	left, right = split2(node, prevKey)
	middle, right = split2(right, key)
	return
}

func (t Treap) GetRandom() (a, b, c, value int) {
	node := t.root
	if node == nil {
		return
	}
	for {
		leftSize := 0
		if node.left != nil {
			leftSize = node.left.size
		}
		next := rand.Intn(node.size)
		if next < leftSize {
			node = node.left
		} else if next == leftSize {
			break
		} else {
			node = node.right
		}
	}
	return node.key.a, node.key.b, node.key.c, node.value
}

func (t Treap) Get(a, b, c int) int {
	key := key3{a, b, c}
	node := t.root
	for node != nil && node.key != key {
		if cmp(node.key, key) < 0 {
			node = node.right
		} else {
			node = node.left
		}
	}
	if node == nil {
		return 0
	}
	return node.value
}

func (t Treap) Set(a, b, c, value int) Treap {
	key := key3{a, b, c}
	left, middle, right := split3(t.root, key)
	if middle == nil {
		middle = newNode(key, value)
	} else {
		middle = middle.copy()
		middle.value = value
	}
	return Treap{merge(left, merge(middle, right))}
}

func (t Treap) Remove(a, b, c int) Treap {
	key := key3{a, b, c}
	left, _, right := split3(t.root, key)
	return Treap{merge(left, right)}
}

func (t Treap) GetBounds3(a, b int) (minC int, maxC int) {
	const inf = 1e9
	prevKey := key3{a, b - 1, +inf}
	nextKey := key3{a, b + 1, -inf}
	middle, _ := split2(t.root, nextKey)
	_, middle = split2(middle, prevKey)
	if middle == nil {
		return
	}
	return middle.minC, middle.maxC
}

func cmp(x, y key3) int {
	if x.a != y.a {
		return x.a - y.a
	}
	if x.b != y.b {
		return x.b - y.b
	}
	if x.c != y.c {
		return x.c - y.c
	}
	return 0
}
