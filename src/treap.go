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
	value    int
	priority int
	left     *node
	right    *node
}

func newNode(key key3, value int) *node {
	return &node{
		key:      key,
		value:    value,
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
	return
}

func split2(node *node, key key3) (left *node, right *node) {
	if node == nil {
		return
	}
	if cmp(node.key, key) <= 0 {
		left = node.copy()
		left.right, right = split2(node.right, key)
	} else {
		right = node.copy()
		left, right.left = split2(node.left, key)
	}
	return
}

func split3(node *node, key key3) (left, middle, right *node) {
	prevKey := key3{key.a, key.b, key.c - 1}
	left, right = split2(node, prevKey)
	middle, right = split2(right, key)
	return
}

func (t Treap) Get(a, b, c int) int {
	_, middle, _ := split3(t.root, key3{a, b, c})
	if middle == nil {
		return 0
	}
	return middle.value
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
