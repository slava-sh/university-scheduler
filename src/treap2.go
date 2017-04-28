package main

import "math/rand"

// Map-like persistent treap.
type Treap2 struct {
	root *node2
}

type key2 struct {
	a int
	b int
}

type node2 struct {
	key      key2
	value    int
	priority int
	left     *node2
	right    *node2
}

func newNode2(key key2, value int) *node2 {
	return &node2{
		key:      key,
		value:    value,
		priority: rand.Int(),
	}
}

func (n *node2) copy() *node2 {
	if n == nil {
		return nil
	}
	copy := new(node2)
	*copy = *n
	return copy
}

func merge2(left, right *node2) (result *node2) {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.priority >= right.priority {
		result = left.copy()
		result.right = merge2(result.right, right)
	} else {
		result = right.copy()
		result.left = merge2(left, result.left)
	}
	return
}

func split2_2(node2 *node2, key key2) (left *node2, right *node2) {
	if node2 == nil {
		return
	}
	if cmp2(node2.key, key) <= 0 {
		left = node2.copy()
		left.right, right = split2_2(node2.right, key)
	} else {
		right = node2.copy()
		left, right.left = split2_2(node2.left, key)
	}
	return
}

func split3_2(node2 *node2, key key2) (left, middle, right *node2) {
	prevKey := key2{key.a, key.b - 1}
	left, right = split2_2(node2, prevKey)
	middle, right = split2_2(right, key)
	return
}

func (t Treap2) Get(a, b int) int {
	key := key2{a, b}
	node2 := t.root
	for node2 != nil && node2.key != key {
		if cmp2(node2.key, key) < 0 {
			node2 = node2.right
		} else {
			node2 = node2.left
		}
	}
	if node2 == nil {
		return 0
	}
	return node2.value
}

func (t Treap2) Set(a, b, value int) Treap2 {
	key := key2{a, b}
	left, middle, right := split3_2(t.root, key)
	if middle == nil {
		middle = newNode2(key, value)
	} else {
		middle = middle.copy()
		middle.value = value
	}
	return Treap2{merge2(left, merge2(middle, right))}
}

func (t Treap2) Inc(a, b int) Treap2 {
	key := key2{a, b}
	left, middle, right := split3_2(t.root, key)
	if middle == nil {
		middle = newNode2(key, 1)
	} else {
		middle = middle.copy()
		middle.value++
	}
	return Treap2{merge2(left, merge2(middle, right))}
}

func (t Treap2) Dec(a, b int) Treap2 {
	key := key2{a, b}
	left, middle, right := split3_2(t.root, key)
	if middle == nil {
		middle = newNode2(key, -1)
	} else if middle.value == 1 {
		middle = nil
	} else {
		middle = middle.copy()
		middle.value--
	}
	return Treap2{merge2(left, merge2(middle, right))}
}

func (t Treap2) Remove(a, b int) Treap2 {
	key := key2{a, b}
	left, _, right := split3_2(t.root, key)
	return Treap2{merge2(left, right)}
}

func cmp2(x, y key2) int {
	if x.a != y.a {
		return x.a - y.a
	}
	if x.b != y.b {
		return x.b - y.b
	}
	return 0
}
