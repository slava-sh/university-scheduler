package main

import "math/rand"

// Map-like persistent treap.
type Treap3 struct {
	root *node3
}

type key3 struct {
	a int
	b int
	c int
}

type node3 struct {
	key      key3
	value    int
	priority int
	left     *node3
	right    *node3
}

func newNode3(key key3, value int) *node3 {
	return &node3{
		key:      key,
		value:    value,
		priority: rand.Int(),
	}
}

func (n *node3) copy() *node3 {
	if n == nil {
		return nil
	}
	copy := new(node3)
	*copy = *n
	return copy
}

func merge3(left, right *node3) (result *node3) {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.priority >= right.priority {
		result = left.copy()
		result.right = merge3(result.right, right)
	} else {
		result = right.copy()
		result.left = merge3(left, result.left)
	}
	return
}

func split2_3(node3 *node3, key key3) (left *node3, right *node3) {
	if node3 == nil {
		return
	}
	if cmp3(node3.key, key) <= 0 {
		left = node3.copy()
		left.right, right = split2_3(node3.right, key)
	} else {
		right = node3.copy()
		left, right.left = split2_3(node3.left, key)
	}
	return
}

func split3_3(node3 *node3, key key3) (left, middle, right *node3) {
	prevKey := key3{key.a, key.b, key.c - 1}
	left, right = split2_3(node3, prevKey)
	middle, right = split2_3(right, key)
	return
}

func (t Treap3) Get(a, b, c int) int {
	key := key3{a, b, c}
	node3 := t.root
	for node3 != nil && node3.key != key {
		if cmp3(node3.key, key) < 0 {
			node3 = node3.right
		} else {
			node3 = node3.left
		}
	}
	if node3 == nil {
		return 0
	}
	return node3.value
}

func (t Treap3) Set(a, b, c, value int) Treap3 {
	key := key3{a, b, c}
	left, middle, right := split3_3(t.root, key)
	if middle == nil {
		middle = newNode3(key, value)
	} else {
		middle = middle.copy()
		middle.value = value
	}
	return Treap3{merge3(left, merge3(middle, right))}
}

func (t Treap3) Remove(a, b, c int) Treap3 {
	key := key3{a, b, c}
	left, _, right := split3_3(t.root, key)
	return Treap3{merge3(left, right)}
}

func cmp3(x, y key3) int {
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
